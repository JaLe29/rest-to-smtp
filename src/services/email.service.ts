import nodemailer, { Transporter } from 'nodemailer';
import { EmailRequest, EmailResponse } from '../types/email';

export class EmailService {
  private createTransporter(config: EmailRequest): Transporter {
    return nodemailer.createTransport({
      host: config.smtp_host,
      port: config.smtp_port,
      secure: config.smtp_port === 465, // true for 465, false for other ports
      auth: {
        user: config.username,
        pass: config.password
      },
      tls: {
        rejectUnauthorized: false // Allow self-signed certificates
      }
    });
  }

  async sendEmail(config: EmailRequest): Promise<EmailResponse> {
    const startTime = Date.now();

    try {
      console.log(`[${new Date().toISOString()}] Email send attempt:`, {
        to: config.to,
        subject: config.subject,
        smtp_host: config.smtp_host,
        smtp_port: config.smtp_port,
        username: config.username
      });

      const transporter = this.createTransporter(config);

      // Verify connection
      console.log(`[${new Date().toISOString()}] Testing SMTP connection...`);
      await transporter.verify();
      console.log(`[${new Date().toISOString()}] SMTP connection verified`);

      // Send email
      const mailOptions = {
        from: config.username,
        to: config.to,
        subject: config.subject,
        text: config.body,
        headers: {
          'X-Mailer': 'rest-to-smtp/1.0',
          'X-Priority': '3'
        }
      };

      console.log(`[${new Date().toISOString()}] Sending email...`);
      const info = await transporter.sendMail(mailOptions);

      const duration = Date.now() - startTime;
      console.log(`[${new Date().toISOString()}] Email sent successfully:`, {
        to: config.to,
        subject: config.subject,
        messageId: info.messageId,
        duration: `${duration}ms`
      });

      return {
        success: true,
        message: 'Email sent successfully',
        messageId: info.messageId,
        duration
      };

    } catch (error) {
      const duration = Date.now() - startTime;
      const errorMessage = error instanceof Error ? error.message : 'Unknown error';

      console.error(`[${new Date().toISOString()}] Email send failed:`, {
        to: config.to,
        subject: config.subject,
        error: errorMessage,
        duration: `${duration}ms`
      });

      return {
        success: false,
        message: `Failed to send email: ${errorMessage}`,
        duration
      };
    }
  }
}
