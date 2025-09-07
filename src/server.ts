import Fastify, { FastifyInstance } from 'fastify';
import cors from '@fastify/cors';
import { EmailController } from './controllers/email.controller';
import { HealthController } from './controllers/health.controller';

const PORT = process.env.PORT || 8080;

async function buildServer(): Promise<FastifyInstance> {
  const fastify = Fastify({
    logger: {
      level: 'info'
    }
  });

  // Register CORS
  await fastify.register(cors, {
    origin: true
  });

  // Initialize controllers
  const emailController = new EmailController();
  const healthController = new HealthController();

  // Routes
  fastify.post('/send-email', {
    handler: emailController.sendEmail.bind(emailController)
  });

  fastify.get('/health', {
    handler: healthController.health.bind(healthController)
  });

  return fastify;
}

async function start() {
  try {
    const server = await buildServer();

    await server.listen({
      port: Number(PORT),
      host: '0.0.0.0'
    });

    console.log(`[${new Date().toISOString()}] REST-to-SMTP server running on port ${PORT}`);
    console.log(`[${new Date().toISOString()}] Available endpoints:`);
    console.log(`[${new Date().toISOString()}]   POST /send-email - Send email via SMTP`);
    console.log(`[${new Date().toISOString()}]   GET  /health     - Health check`);

  } catch (err) {
    console.error(`[${new Date().toISOString()}] Server failed to start:`, err);
    process.exit(1);
  }
}

// Handle graceful shutdown
process.on('SIGINT', async () => {
  console.log(`[${new Date().toISOString()}] Received SIGINT, shutting down gracefully...`);
  process.exit(0);
});

process.on('SIGTERM', async () => {
  console.log(`[${new Date().toISOString()}] Received SIGTERM, shutting down gracefully...`);
  process.exit(0);
});

start();
