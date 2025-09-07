# REST-to-SMTP

A simple Go application that provides a REST API for sending emails via SMTP.

## Features

- REST API endpoint for sending emails
- SMTP configuration via API request
- JSON request/response format
- CORS support
- Health check endpoint
- Input validation
- Docker containerization
- GitHub Actions CI/CD pipeline
- Multi-stage Docker build for optimal image size

## Project Structure

```
rest-to-smtp/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   │   ├── email.go
│   │   └── health.go
│   ├── models/              # Data structures
│   │   └── email.go
│   └── services/            # Business logic
│       └── email.go
├── Dockerfile               # Multi-stage Docker build
├── .github/workflows/      # GitHub Actions CI/CD pipeline
├── go.mod                  # Go module definition
└── README.md
```

## Installation

### Option 1: Using Docker (Recommended)

1. Pull the Docker image from Docker Hub:
```bash
docker pull your-github-username/rest-to-smtp:latest
```

2. Run the container:
```bash
docker run -p 8080:8080 your-github-username/rest-to-smtp:latest
```

> **Note:** The image name follows the pattern `{github-username}/{repository-name}` and is automatically built and pushed by GitHub Actions.

### Option 2: Building from Source

1. Make sure you have Go installed (version 1.21 or higher)
2. Clone or download this repository
3. Navigate to the project directory
4. Run the application:

```bash
go run ./cmd/server
```

The server will start on port 8080.

### Option 3: Building Docker Image Locally

1. Clone the repository
2. Build the Docker image:
```bash
docker build -t rest-to-smtp .
```

3. Run the container:
```bash
docker run -p 8080:8080 rest-to-smtp
```

## API Endpoints

### POST /send-email

Send an email via SMTP.

**Request Body:**
```json
{
  "smtp_host": "smtp.gmail.com",
  "smtp_port": "587",
  "username": "your-email@gmail.com",
  "password": "your-password",
  "to": "recipient@example.com",
  "subject": "Test Email",
  "body": "This is a test email body"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Email sent successfully"
}
```

### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "service": "rest-to-smtp"
}
```

## Usage Examples

### Using curl

```bash
curl -X POST http://localhost:8080/send-email \
  -H "Content-Type: application/json" \
  -d '{
    "smtp_host": "smtp.gmail.com",
    "smtp_port": "587",
    "username": "your-email@gmail.com",
    "password": "your-app-password",
    "to": "recipient@example.com",
    "subject": "Test Email",
    "body": "Hello from REST-to-SMTP!"
  }'
```

### Using JavaScript (fetch)

```javascript
const response = await fetch('http://localhost:8080/send-email', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    smtp_host: 'smtp.gmail.com',
    smtp_port: '587',
    username: 'your-email@gmail.com',
    password: 'your-app-password',
    to: 'recipient@example.com',
    subject: 'Test Email',
    body: 'Hello from REST-to-SMTP!'
  })
});

const result = await response.json();
console.log(result);
```

## SMTP Configuration Examples

### Gmail
- Host: `smtp.gmail.com`
- Port: `587` (TLS) or `465` (SSL)
- Username: Your Gmail address
- Password: App password (not your regular password)

### Outlook/Hotmail
- Host: `smtp-mail.outlook.com`
- Port: `587`
- Username: Your Outlook email
- Password: Your Outlook password

### Custom SMTP Server
- Host: Your SMTP server address
- Port: Usually `587` or `465`
- Username: Your SMTP username
- Password: Your SMTP password

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `400 Bad Request`: Invalid JSON, missing required fields, or SMTP connection failed
- `405 Method Not Allowed`: Wrong HTTP method
- `500 Internal Server Error`: SMTP sending failed

### Timeout Protection

The application includes built-in timeout protection:

- **SMTP Connection Test**: 10 seconds timeout for initial connection test
- **Email Sending**: 30 seconds timeout for complete email sending process
- **Fast Failure**: Invalid SMTP servers are detected quickly without hanging

### Common Error Messages

- `"SMTP connection failed: cannot connect to SMTP server"` - Server unreachable
- `"SMTP connection timeout after 30 seconds"` - Operation took too long
- `"Username and Password not accepted"` - Authentication failed
- `"invalid SMTP port. Use 25, 587, or 465"` - Invalid port number

## Security Notes

- This application is designed for development and testing purposes
- In production, consider implementing authentication and rate limiting
- Never hardcode SMTP credentials in your client applications
- Use environment variables or secure configuration management for sensitive data

## GitHub Actions CI/CD

This project includes a GitHub Actions CI/CD pipeline that automatically:

1. **Tests** - Runs Go tests and linting on pull requests and pushes
2. **Builds** - Creates Docker images for different branches
3. **Deploys** - Pushes images to Docker Hub

### Setup for GitHub Actions

1. **Create Docker Hub Access Token:**
   - Go to [hub.docker.com](https://hub.docker.com) and sign in
   - Click your username → **Account Settings**
   - Go to **Security** → **New Access Token**
   - Name: "github-actions", Permissions: **Read, Write, Delete**
   - Click **Generate** and **copy the token** (shown only once!)

2. **Set up Docker Hub credentials in GitHub:**
   - Go to your GitHub repository → **Settings** → **Secrets and variables** → **Actions**
   - Click **New repository secret** and add:
     - `DOCKER_USERNAME`: Your Docker Hub username
     - `DOCKER_PASSWORD`: Your Docker Hub access token (NOT your password!)

3. **Workflow triggers:**
   - `push`: Runs on main and develop branches
   - `pull_request`: Runs on PRs targeting main or develop branches

4. **Pipeline stages:**
   - `test`: Runs Go tests, linting, and dependency verification
   - `build-and-push`: Builds and pushes Docker images (only on push events)
   - `deploy`: Deploys to production (only on main branch)

5. **Image tags:**
   - `latest`: For main branch
   - `develop`: For develop branch
   - `{branch}-{commit-sha}`: For other branches

### Manual Deployment

To manually deploy to Docker Hub:

```bash
# Build and tag the image
docker build -t your-github-username/rest-to-smtp:latest .

# Push to Docker Hub
docker push your-github-username/rest-to-smtp:latest
```

### Workflow Features

- **Multi-platform builds**: Supports both AMD64 and ARM64 architectures
- **Caching**: Uses GitHub Actions cache for faster builds
- **Security**: Uses Docker Buildx for secure builds
- **Automatic tagging**: Creates appropriate tags based on branch and commit
- **Environment protection**: Production deployment requires manual approval

## Building for Production

### Using Go

To build a binary for production:

```bash
go build -o rest-to-smtp ./cmd/server
```

Then run the binary:

```bash
./rest-to-smtp
```

### Using Docker

The Dockerfile uses multi-stage build for optimal image size:

```bash
# Build the image
docker build -t rest-to-smtp .

# Run the container
docker run -p 8080:8080 rest-to-smtp
```

The final image is based on Alpine Linux and includes:
- Non-root user for security
- Health check endpoint
- Minimal attack surface
- Small image size (~15MB)
