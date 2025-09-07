export interface EmailRequest {
  smtp_host: string;
  smtp_port: number;
  username: string;
  password: string;
  to: string;
  subject: string;
  body: string;
}

export interface EmailResponse {
  success: boolean;
  message: string;
  messageId?: string;
  duration?: number;
}

export interface HealthResponse {
  status: string;
  timestamp: string;
  version: string;
}
