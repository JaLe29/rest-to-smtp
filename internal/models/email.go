package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// EmailRequest represents the structure of the incoming email request
type EmailRequest struct {
	SMTPHost string `json:"smtp_host"`
	SMTPPort int    `json:"smtp_port"`
	Username string `json:"username"`
	Password string `json:"password"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

// UnmarshalJSON custom unmarshaling for SMTPPort to handle both string and int
func (e *EmailRequest) UnmarshalJSON(data []byte) error {
	type Alias EmailRequest
	aux := &struct {
		SMTPPort interface{} `json:"smtp_port"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle SMTPPort conversion
	switch v := aux.SMTPPort.(type) {
	case string:
		port, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid smtp_port: %v", err)
		}
		e.SMTPPort = port
	case float64:
		e.SMTPPort = int(v)
	case int:
		e.SMTPPort = v
	default:
		return fmt.Errorf("smtp_port must be a number")
	}

	return nil
}

// EmailResponse represents the response structure
type EmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
