package models

// EmailRequest represents the structure of the incoming email request
type EmailRequest struct {
	SMTPHost string `json:"smtp_host"`
	SMTPPort string `json:"smtp_port"`
	Username string `json:"username"`
	Password string `json:"password"`
	To       string `json:"to"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
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
