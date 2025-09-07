# REST-to-SMTP

A modern REST API for sending emails via SMTP, built with Node.js, TypeScript, Fastify, and Nodemailer.

## Features

- 🚀 **Fast & Modern**: Built with Fastify (faster than Express)
- 📧 **SMTP Support**: Works with any SMTP server (Gmail, Outlook, custom servers)
- 🔒 **Type Safety**: Full TypeScript support with strict typing
- ✅ **Validation**: Zod schema validation for all inputs
- 🏗️ **Layered Architecture**: Clean separation of concerns
- 🐳 **Docker Ready**: Multi-stage Docker build with security best practices
- 🔄 **CI/CD**: GitHub Actions pipeline for automated testing and deployment
- 📊 **Structured Logging**: Detailed logging for monitoring and debugging
- 🛡️ **Security**: Non-root user, health checks, and input validation

## Architecture

```
src/
├── controllers/     # HTTP request handlers
├── services/        # Business logic layer
├── validators/      # Zod validation schemas
├── types/          # TypeScript type definitions
├── utils/          # Utility functions
└── server.ts       # Application entry point
```

## Quick Start

### Prerequisites

- Node.js 18+
- npm or yarn
- Docker (optional)

### Local Development

1. **Clone and install dependencies:**
   ```bash
   git clone <repository-url>
   cd rest-to-smtp
   npm install
   ```

2. **Build the application:**
   ```bash
   npm run build
   ```

3. **Start the server:**
   ```bash
   npm start
   # or for development with hot reload:
   npm run dev
   ```

4. **Test the API:**
   ```bash
   # Health check
   curl http://localhost:8080/health

   # Send email
   curl -X POST http://localhost:8080/send-email \
     -H "Content-Type: application/json" \
     -d '{
       "smtp_host": "smtp.gmail.com",
       "smtp_port": 587,
       "username": "your-email@gmail.com",
       "password": "your-app-password",
       "to": "recipient@example.com",
       "subject": "Test Email",
       "body": "Hello from REST-to-SMTP!"
     }'
   ```

### Docker

1. **Build the image:**
   ```bash
   docker build -t rest-to-smtp .
   ```

2. **Run the container:**
   ```bash
   docker run -d -p 8080:8080 --name rest-to-smtp rest-to-smtp
   ```

3. **Test the container:**
   ```bash
   curl http://localhost:8080/health
   ```

## API Endpoints

### POST /send-email

Send an email via SMTP.

**Request Body:**
```json
{
  "smtp_host": "smtp.gmail.com",
  "smtp_port": 587,
  "username": "your-email@gmail.com",
  "password": "your-password",
  "to": "recipient@example.com",
  "subject": "Email Subject",
  "body": "Email content"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Email sent successfully",
  "messageId": "<message-id>",
  "duration": 1234
}
```

### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2025-09-07T14:10:09.064Z",
  "version": "1.0.0"
}
```

## SMTP Configuration

### Supported Ports
- **25**: Plain SMTP (no encryption)
- **587**: SMTP with STARTTLS (recommended)
- **465**: SMTPS (implicit SSL)

### Popular SMTP Providers

#### Gmail
```json
{
  "smtp_host": "smtp.gmail.com",
  "smtp_port": 587,
  "username": "your-email@gmail.com",
  "password": "your-app-password"
}
```

#### Outlook/Hotmail
```json
{
  "smtp_host": "smtp-mail.outlook.com",
  "smtp_port": 587,
  "username": "your-email@outlook.com",
  "password": "your-password"
}
```

#### Custom SMTP Server
```json
{
  "smtp_host": "mail.yourdomain.com",
  "smtp_port": 587,
  "username": "your-email@yourdomain.com",
  "password": "your-password"
}
```

## Validation

The API validates all inputs using Zod schemas:

- **SMTP Host**: Must be a valid domain name
- **Port**: Must be 25, 587, or 465
- **Username**: Must be a valid email address
- **Password**: Required (minimum 1 character)
- **To**: Must be a valid email address
- **Subject**: Required, max 200 characters
- **Body**: Required, max 10,000 characters

## CI/CD Pipeline

The project includes GitHub Actions workflows:

### Test Workflow (`.github/workflows/test.yml`)
- Runs on every push and pull request
- Tests with Node.js 18 and 20
- TypeScript compilation check
- Security audit
- Linting (if configured)

### CI/CD Workflow (`.github/workflows/ci-cd.yml`)
- Runs on push to main branch
- Builds and pushes Docker image to GitHub Container Registry
- Multi-platform build (linux/amd64, linux/arm64)
- Automated deployment notifications

## Docker Image

The Docker image is automatically built and pushed to GitHub Container Registry:

```bash
# Pull the latest image
docker pull ghcr.io/your-username/rest-to-smtp:latest

# Run the container
docker run -d -p 8080:8080 ghcr.io/your-username/rest-to-smtp:latest
```

## Environment Variables

- `PORT`: Server port (default: 8080)

## Security Features

- Non-root user in Docker container
- Input validation and sanitization
- Structured logging for monitoring
- Health checks for container orchestration
- CORS support for web applications

## Error Handling

The API provides detailed error messages:

```json
{
  "success": false,
  "message": "Validation failed: smtp_port: Unsupported SMTP port. Supported ports: 25, 587, 465"
}
```

Common error scenarios:
- Invalid SMTP configuration
- Network connectivity issues
- Authentication failures
- Validation errors

## Development

### Project Structure
```
rest-to-smtp/
├── src/
│   ├── controllers/     # HTTP handlers
│   ├── services/        # Business logic
│   ├── validators/      # Zod schemas
│   ├── types/          # TypeScript types
│   ├── utils/          # Utilities
│   └── server.ts       # Entry point
├── dist/               # Compiled JavaScript
├── .github/workflows/  # CI/CD pipelines
├── Dockerfile         # Docker configuration
├── package.json       # Dependencies
└── tsconfig.json      # TypeScript config
```

### Scripts
- `npm run build`: Compile TypeScript
- `npm start`: Start production server
- `npm run dev`: Start development server with hot reload
- `npm test`: Run tests (if available)

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For issues and questions, please open an issue on GitHub.