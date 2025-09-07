import { z } from 'zod';

export const emailRequestSchema = z.object({
  smtp_host: z.string()
    .min(1, 'SMTP host is required')
    .regex(/^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/, 'Invalid SMTP host format'),

  smtp_port: z.number()
    .int('Port must be an integer')
    .min(1, 'Port must be greater than 0')
    .max(65535, 'Port must be less than 65536')
    .refine((port) => [25, 587, 465].includes(port), {
      message: 'Unsupported SMTP port. Supported ports: 25, 587, 465'
    }),

  username: z.string()
    .min(1, 'Username is required')
    .email('Invalid email address format'),

  password: z.string()
    .min(1, 'Password is required'),

  to: z.string()
    .min(1, 'Recipient email is required')
    .email('Invalid recipient email address'),

  subject: z.string()
    .min(1, 'Subject is required')
    .max(200, 'Subject must be less than 200 characters'),

  body: z.string()
    .min(1, 'Email body is required')
    .max(10000, 'Email body must be less than 10000 characters')
});

export type EmailRequestType = z.infer<typeof emailRequestSchema>;
