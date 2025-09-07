import { FastifyRequest, FastifyReply } from 'fastify';
import { HealthResponse } from '../types/email';

export class HealthController {
  async health(request: FastifyRequest, reply: FastifyReply): Promise<void> {
    const response: HealthResponse = {
      status: 'healthy',
      timestamp: new Date().toISOString(),
      version: '1.0.0'
    };

    reply.status(200).send(response);
  }
}
