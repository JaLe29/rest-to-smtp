import { FastifyRequest, FastifyReply } from 'fastify';
import { EmailService } from '../services/email.service';
import { emailRequestSchema } from '../validators/email';
import { EmailRequest, EmailResponse } from '../types/email';

export class EmailController {
  private emailService: EmailService;

  constructor() {
    this.emailService = new EmailService();
  }

  async sendEmail(request: FastifyRequest, reply: FastifyReply): Promise<void> {
    const startTime = Date.now();
    const clientIP = request.ip;

    try {
      // Log incoming request
      console.log(`[${new Date().toISOString()}] Email send request received:`, {
        client_ip: clientIP,
        user_agent: request.headers['user-agent']
      });

      // Validate request body
      const validationResult = emailRequestSchema.safeParse(request.body);

      if (!validationResult.success) {
        const errors = validationResult.error.errors.map(err =>
          `${err.path.join('.')}: ${err.message}`
        ).join(', ');

        console.error(`[${new Date().toISOString()}] Validation failed:`, {
          client_ip: clientIP,
          errors
        });

        const response: EmailResponse = {
          success: false,
          message: `Validation failed: ${errors}`
        };

        reply.status(400).send(response);
        return;
      }

      const emailConfig: EmailRequest = validationResult.data;

      // Log validated request
      console.log(`[${new Date().toISOString()}] Email request validated:`, {
        client_ip: clientIP,
        to: emailConfig.to,
        subject: emailConfig.subject,
        smtp_host: emailConfig.smtp_host,
        smtp_port: emailConfig.smtp_port,
        username: emailConfig.username
      });

      // Send email
      const result = await this.emailService.sendEmail(emailConfig);

      // Log completion
      const duration = Date.now() - startTime;
      if (result.success) {
        console.log(`[${new Date().toISOString()}] Email request completed successfully:`, {
          client_ip: clientIP,
          to: emailConfig.to,
          subject: emailConfig.subject,
          duration: `${duration}ms`
        });
        reply.status(200).send(result);
      } else {
        console.error(`[${new Date().toISOString()}] Email request failed:`, {
          client_ip: clientIP,
          to: emailConfig.to,
          subject: emailConfig.subject,
          error: result.message,
          duration: `${duration}ms`
        });
        reply.status(500).send(result);
      }

    } catch (error) {
      const duration = Date.now() - startTime;
      const errorMessage = error instanceof Error ? error.message : 'Unknown error';

      console.error(`[${new Date().toISOString()}] Unexpected error in email controller:`, {
        client_ip: clientIP,
        error: errorMessage,
        duration: `${duration}ms`
      });

      const response: EmailResponse = {
        success: false,
        message: `Internal server error: ${errorMessage}`,
        duration
      };

      reply.status(500).send(response);
    }
  }
}
