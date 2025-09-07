package services

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/wneessen/go-mail"
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

	// Create new message
	message := mail.NewMsg()

	// Set From header
	if err := message.From(req.Username); err != nil {
		return fmt.Errorf("failed to set From address: %v", err)
	}

	// Set To header
	if err := message.To(req.To); err != nil {
		return fmt.Errorf("failed to set To address: %v", err)
	}

	// Set subject and body
	message.Subject(req.Subject)
	message.SetBodyString(mail.TypeTextPlain, req.Body)

	// Create SMTP client
	client, err := mail.NewClient(req.SMTPHost,
		mail.WithPort(req.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(req.Username),
		mail.WithPassword(req.Password),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %v", err)
	}

	// Send email
	slog.Info("Sending email", "to", req.To, "subject", req.Subject)
	if err := client.DialAndSend(message); err != nil {
		duration := time.Since(startTime)
		slog.Error("Email send failed",
			"to", req.To,
			"subject", req.Subject,
			"smtp_host", req.SMTPHost,
			"smtp_port", req.SMTPPort,
			"error", err.Error(),
			"duration_ms", duration.Milliseconds(),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		return fmt.Errorf("failed to send mail: %v", err)
	}

	// Log successful email
	duration := time.Since(startTime)
	slog.Info("Email sent successfully",
		"to", req.To,
		"subject", req.Subject,
		"smtp_host", req.SMTPHost,
		"smtp_port", req.SMTPPort,
		"duration_ms", duration.Milliseconds(),
		"timestamp", time.Now().Format(time.RFC3339),
	)

	return nil
}



// ValidateEmailRequest validates the email request
func (s *EmailService) ValidateEmailRequest(req models.EmailRequest) error {
	// Check required fields
	if req.SMTPHost == "" {
		return fmt.Errorf("smtp_host is required")
	}
	if req.SMTPPort == 0 {
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
	if req.Body == "" {
		return fmt.Errorf("body is required")
	}

	// Validate SMTP host format
	if !strings.Contains(req.SMTPHost, ".") {
		return fmt.Errorf("invalid SMTP host format")
	}

	// Validate port
	allowedPorts := []int{25, 587, 465}
	portValid := false
	for _, port := range allowedPorts {
		if req.SMTPPort == port {
			portValid = true
			break
		}
	}
	if !portValid {
		return fmt.Errorf("unsupported SMTP port: %d (supported: 25, 587, 465)", req.SMTPPort)
	}

	// Validate email addresses
	if !strings.Contains(req.To, "@") {
		return fmt.Errorf("invalid recipient email address")
	}
	if !strings.Contains(req.Username, "@") {
		return fmt.Errorf("invalid sender email address")
	}

	return nil
}