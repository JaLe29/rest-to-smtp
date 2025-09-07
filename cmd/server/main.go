package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"rest-to-smtp/internal/handlers"
	"rest-to-smtp/internal/services"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Log application startup
	slog.Info("Starting REST-to-SMTP server",
		"version", "1.0.0",
		"gin_mode", gin.Mode(),
	)

	// Initialize services
	emailService := services.NewEmailService()

	// Initialize handlers
	emailHandler := handlers.NewEmailHandler(emailService)
	healthHandler := handlers.NewHealthHandler()

	// Create Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Define routes
	r.POST("/send-email", emailHandler.SendEmail)
	r.GET("/health", healthHandler.Health)

	// Start server
	port := "8080"
	log.Printf("Starting REST-to-SMTP server on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  POST /send-email - Send email via SMTP")
	log.Printf("  GET  /health     - Health check")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
