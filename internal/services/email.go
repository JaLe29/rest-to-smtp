package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"strings"
	"time"

	"rest-to-smtp/internal/models"
)

// EmailService handles email sending operations
type EmailService struct{}

// NewEmailService creates a new email service instance
func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendEmail sends an email using the provided SMTP configuration
func (s *EmailService) SendEmail(req models.EmailRequest) error {
	startTime := time.Now()

	// Log email attempt
	slog.Info("Email send attempt started",
		"to", req.To,
		"subject", req.Subject,
		"smtp_host", req.SMTPHost,
		"smtp_port", req.SMTPPort,
		"username", req.Username,
		"timestamp", startTime.Format(time.RFC3339),
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Email content
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", req.To, req.Subject, req.Body)

	// Use a channel to handle the SMTP operation with timeout
	errChan := make(chan error, 1)

	go func() {
		var err error

		// Choose connection method based on port
		switch req.SMTPPort {
		case "465":
			// Port 465: SMTPS (SMTP over SSL)
			slog.Info("Using SMTPS connection (port 465)", "host", req.SMTPHost)
			err = s.sendEmailSMTPS(req, msg)
		case "587":
			// Port 587: SMTP with STARTTLS
			slog.Info("Using STARTTLS connection (port 587)", "host", req.SMTPHost)
			err = s.sendEmailSTARTTLS(req, msg)
		case "25":
			// Port 25: Plain SMTP
			slog.Info("Using plain SMTP connection (port 25)", "host", req.SMTPHost)
			err = s.sendEmailPlain(req, msg)
		default:
			err = fmt.Errorf("unsupported port: %s", req.SMTPPort)
		}

		errChan <- err
	}()

	select {
	case err := <-errChan:
		duration := time.Since(startTime)
		if err != nil {
			// Log failed email
			slog.Error("Email send failed",
				"to", req.To,
				"subject", req.Subject,
				"smtp_host", req.SMTPHost,
				"smtp_port", req.SMTPPort,
				"error", err.Error(),
				"duration_ms", duration.Milliseconds(),
				"timestamp", time.Now().Format(time.RFC3339),
			)
		} else {
			// Log successful email
			slog.Info("Email sent successfully",
				"to", req.To,
				"subject", req.Subject,
				"smtp_host", req.SMTPHost,
				"smtp_port", req.SMTPPort,
				"duration_ms", duration.Milliseconds(),
				"timestamp", time.Now().Format(time.RFC3339),
			)
		}
		return err
	case <-ctx.Done():
		duration := time.Since(startTime)
		// Log timeout
		slog.Error("Email send timeout",
			"to", req.To,
			"subject", req.Subject,
			"smtp_host", req.SMTPHost,
			"smtp_port", req.SMTPPort,
			"timeout_seconds", 30,
			"duration_ms", duration.Milliseconds(),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		return fmt.Errorf("SMTP connection timeout after 30 seconds")
	}
}

// sendEmailSMTPS sends email using SMTPS (port 465)
func (s *EmailService) sendEmailSMTPS(req models.EmailRequest, msg string) error {
	addr := fmt.Sprintf("%s:%s", req.SMTPHost, req.SMTPPort)

	// Connect to server with SSL
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: req.SMTPHost,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to SMTPS server: %v", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, req.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// Authenticate
	auth := smtp.PlainAuth("", req.Username, req.Password, req.SMTPHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Send email
	if err := client.Mail(req.Username); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	if err := client.Rcpt(req.To); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	if _, err := writer.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	return nil
}

// sendEmailSTARTTLS sends email using SMTP with STARTTLS (port 587)
func (s *EmailService) sendEmailSTARTTLS(req models.EmailRequest, msg string) error {
	addr := fmt.Sprintf("%s:%s", req.SMTPHost, req.SMTPPort)

	// Connect to server
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, req.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// Start TLS
	if err := client.StartTLS(&tls.Config{ServerName: req.SMTPHost}); err != nil {
		return fmt.Errorf("failed to start TLS: %v", err)
	}

	// Authenticate
	auth := smtp.PlainAuth("", req.Username, req.Password, req.SMTPHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %v", err)
	}

	// Send email
	if err := client.Mail(req.Username); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	if err := client.Rcpt(req.To); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %v", err)
	}

	if _, err := writer.Write([]byte(msg)); err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	return nil
}

// sendEmailPlain sends email using plain SMTP (port 25)
func (s *EmailService) sendEmailPlain(req models.EmailRequest, msg string) error {
	addr := fmt.Sprintf("%s:%s", req.SMTPHost, req.SMTPPort)
	auth := smtp.PlainAuth("", req.Username, req.Password, req.SMTPHost)

	return smtp.SendMail(addr, auth, req.Username, []string{req.To}, []byte(msg))
}

// ValidateEmailRequest validates the email request
func (s *EmailService) ValidateEmailRequest(req models.EmailRequest) error {
	// Check required fields
	if req.SMTPHost == "" {
		return fmt.Errorf("smtp_host is required")
	}
	if req.SMTPPort == "" {
		return fmt.Errorf("smtp_port is required")
	}
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if req.To == "" {
		return fmt.Errorf("to is required")
	}
	if req.Subject == "" {
		return fmt.Errorf("subject is required")
	}

	// Basic email validation
	if !strings.Contains(req.To, "@") {
		return fmt.Errorf("invalid email address format")
	}

	// Validate SMTP host format
	if !strings.Contains(req.SMTPHost, ".") {
		return fmt.Errorf("invalid SMTP host format")
	}

	// Validate port (should be numeric)
	if req.SMTPPort != "25" && req.SMTPPort != "587" && req.SMTPPort != "465" {
		return fmt.Errorf("invalid SMTP port. Use 25, 587, or 465")
	}

	return nil
}

// TestSMTPConnection tests if SMTP server is reachable
func (s *EmailService) TestSMTPConnection(req models.EmailRequest) error {
	startTime := time.Now()
	addr := fmt.Sprintf("%s:%s", req.SMTPHost, req.SMTPPort)

	// Log connection test attempt
	slog.Info("SMTP connection test started",
		"smtp_host", req.SMTPHost,
		"smtp_port", req.SMTPPort,
		"address", addr,
		"timestamp", startTime.Format(time.RFC3339),
	)

	// Test connection based on port
	switch req.SMTPPort {
	case "465":
		// Test SSL connection for port 465
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, &tls.Config{
			ServerName: req.SMTPHost,
		})
		if err != nil {
			duration := time.Since(startTime)
			slog.Error("SMTPS connection test failed",
				"smtp_host", req.SMTPHost,
				"smtp_port", req.SMTPPort,
				"error", err.Error(),
				"duration_ms", duration.Milliseconds(),
				"timestamp", time.Now().Format(time.RFC3339),
			)
			return fmt.Errorf("cannot connect to SMTPS server %s: %v", addr, err)
		}
		defer conn.Close()
		slog.Info("SMTPS connection test successful", "host", req.SMTPHost)

	case "587", "25":
		// Test plain connection for ports 587 and 25
		conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
		if err != nil {
			duration := time.Since(startTime)
			slog.Error("SMTP connection test failed",
				"smtp_host", req.SMTPHost,
				"smtp_port", req.SMTPPort,
				"error", err.Error(),
				"duration_ms", duration.Milliseconds(),
				"timestamp", time.Now().Format(time.RFC3339),
			)
			return fmt.Errorf("cannot connect to SMTP server %s: %v", addr, err)
		}
		defer conn.Close()
		slog.Info("SMTP connection test successful", "host", req.SMTPHost)

	default:
		slog.Error("Unsupported SMTP port", "port", req.SMTPPort)
		return fmt.Errorf("unsupported port: %s", req.SMTPPort)
	}

	duration := time.Since(startTime)
	slog.Info("SMTP connection test completed",
		"smtp_host", req.SMTPHost,
		"smtp_port", req.SMTPPort,
		"duration_ms", duration.Milliseconds(),
		"timestamp", time.Now().Format(time.RFC3339),
	)

	return nil
}
