package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"rest-to-smtp/internal/models"
	"rest-to-smtp/internal/services"
)

// EmailHandler handles email-related HTTP requests
type EmailHandler struct {
	emailService *services.EmailService
}

// NewEmailHandler creates a new email handler instance
func NewEmailHandler(emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

// SendEmail handles the POST request for sending emails
func (h *EmailHandler) SendEmail(c *gin.Context) {
	startTime := time.Now()
	clientIP := c.ClientIP()

	// Log incoming request
	slog.Info("Email send request received",
		"client_ip", clientIP,
		"user_agent", c.GetHeader("User-Agent"),
		"timestamp", startTime.Format(time.RFC3339),
	)

	var emailReq models.EmailRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&emailReq); err != nil {
		slog.Error("Invalid JSON in email request",
			"client_ip", clientIP,
			"error", err.Error(),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		c.JSON(http.StatusBadRequest, models.EmailResponse{
			Success: false,
			Message: "Invalid JSON format: " + err.Error(),
		})
		return
	}

	// Log request details (without password)
	slog.Info("Email request validated",
		"client_ip", clientIP,
		"to", emailReq.To,
		"subject", emailReq.Subject,
		"smtp_host", emailReq.SMTPHost,
		"smtp_port", emailReq.SMTPPort,
		"username", emailReq.Username,
		"timestamp", time.Now().Format(time.RFC3339),
	)

	// Validate request
	if err := h.emailService.ValidateEmailRequest(emailReq); err != nil {
		slog.Error("Email request validation failed",
			"client_ip", clientIP,
			"to", emailReq.To,
			"subject", emailReq.Subject,
			"error", err.Error(),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		c.JSON(http.StatusBadRequest, models.EmailResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Test SMTP connection first
	if err := h.emailService.TestSMTPConnection(emailReq); err != nil {
		slog.Error("SMTP connection test failed in handler",
			"client_ip", clientIP,
			"to", emailReq.To,
			"subject", emailReq.Subject,
			"smtp_host", emailReq.SMTPHost,
			"error", err.Error(),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		c.JSON(http.StatusBadRequest, models.EmailResponse{
			Success: false,
			Message: fmt.Sprintf("SMTP connection failed: %v", err),
		})
		return
	}

	// Send email
	if err := h.emailService.SendEmail(emailReq); err != nil {
		duration := time.Since(startTime)
		slog.Error("Email sending failed in handler",
			"client_ip", clientIP,
			"to", emailReq.To,
			"subject", emailReq.Subject,
			"smtp_host", emailReq.SMTPHost,
			"error", err.Error(),
			"duration_ms", duration.Milliseconds(),
			"timestamp", time.Now().Format(time.RFC3339),
		)
		c.JSON(http.StatusInternalServerError, models.EmailResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send email: %v", err),
		})
		return
	}

	// Success response
	duration := time.Since(startTime)
	slog.Info("Email request completed successfully",
		"client_ip", clientIP,
		"to", emailReq.To,
		"subject", emailReq.Subject,
		"smtp_host", emailReq.SMTPHost,
		"duration_ms", duration.Milliseconds(),
		"timestamp", time.Now().Format(time.RFC3339),
	)

	c.JSON(http.StatusOK, models.EmailResponse{
		Success: true,
		Message: "Email sent successfully",
	})
}
