# Service Template 🚀

A simple, stylish Go-based web service template for homelab automation with Docker containerization, CI/CD, and multi-environment support.

## Features

- ✨ Simple & stylish web interface
- 🐳 Docker containerized
- 🔄 Automated CI/CD with GitHub Actions
- 🌍 Multi-environment support (Production & Staging)
- 💚 Health check endpoints
- 📊 Environment-aware UI (different styles for prod/dev)

## Quick Start

### Local Development

```bash
# Clone the repository
git clone https://github.com/MrCodeEU/service-template.git
cd service-template

# Run with Go
go run main.go

# Or build and run
go build -o service-template
./service-template
```

Visit `http://localhost:8080`

### Docker

#### Production Environment

```bash
# Build and run production
docker-compose up -d

# Access at http://localhost:8080
```

#### Development Environment

```bash
# Build and run dev/staging
docker-compose -f docker-compose.dev.yml up -d

# Access at http://localhost:8081
```

#### Pull from GitHub Container Registry

```bash
# Production (latest tag)
docker pull ghcr.io/mrcodeu/service-template:latest
docker run -p 8080:8080 \
  -e DISPLAY_MESSAGE="Production is live!" \
  -e ENVIRONMENT=production \
  ghcr.io/mrcodeu/service-template:latest

# Development (dev tag)
docker pull ghcr.io/mrcodeu/service-template:dev
docker run -p 8081:8080 \
  -e DISPLAY_MESSAGE="Dev environment testing" \
  -e ENVIRONMENT=development \
  ghcr.io/mrcodeu/service-template:dev
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DISPLAY_MESSAGE` | Custom message to display | `Welcome to Service Template!` |
| `ENVIRONMENT` | Environment name (affects UI styling) | `production` |
| `VERSION` | Service version | `1.0.0` |

## Environment Configurations

### Production (`latest` tag)
- **Message**: "🚀 Production Service Running!"
- **Environment**: production
- **UI Theme**: Green gradient with rocket emoji
- **Trigger**: Push to `main` branch

### Development/Staging (`dev` tag)
- **Message**: "🛠️ Development/Staging Environment - Testing in Progress"
- **Environment**: development
- **UI Theme**: Pink/red gradient with tools emoji
- **Trigger**: Push to any `dev` branch or branch containing `dev`

## Endpoints

- `GET /` - Main web interface
- `GET /health` - Health check (returns JSON)
- `GET /ready` - Readiness check (returns JSON)

## CI/CD Workflow

The GitHub Actions workflow automatically:

1. **Tests** - Runs on all pushes and PRs
   - Go tests with race detection
   - Code coverage upload to Codecov
   - golangci-lint checks

2. **Build & Push** - Runs on pushes (not PRs)
   - Builds Docker image
   - Tags appropriately:
     - `latest` for main branch
     - `dev` for dev branches
     - `sha-<commit>` for all pushes
     - Semantic versions for tags
   - Pushes to GitHub Container Registry

3. **Deploy** - Triggers after successful build
   - **Main branch** → Triggers production deployment with `latest` tag
   - **Dev branches** → Triggers staging deployment with `dev` tag
   - Sends webhook to `homelab-automation` repository

### Required Secrets

Create these secrets in your GitHub repository:

- `GITHUB_TOKEN` - Automatically provided by GitHub
- `DISPATCH_TOKEN` - Personal access token with `repo` scope for triggering deployments
  - Create at: https://github.com/settings/tokens
  - Scopes needed: `repo`

## Project Structure

```
.
├── .github/
│   └── workflows/
│       └── build-and-deploy.yml    # Combined CI/CD workflow
├── templates/
│   └── index.html                  # Stylish web interface
├── main.go                         # Go web server
├── Dockerfile                      # Multi-stage Docker build
├── docker-compose.yml              # Production compose file
├── docker-compose.dev.yml          # Development compose file
├── .env.example                    # Production env example
├── .env.dev.example                # Development env example
├── .golangci.yml                   # Linter configuration
├── go.mod                          # Go module definition
└── README.md                       # This file
```

## Homelab Integration

This template is designed to integrate with your homelab automation system. When code is pushed:

1. GitHub Actions builds and pushes a Docker image
2. A webhook is sent to `homelab-automation` repository
3. The automation system pulls the appropriate image tag (`latest` or `dev`)
4. Service is deployed to the corresponding environment

### Webhook Payload

```json
{
  "event_type": "service-update",
  "client_payload": {
    "service": "service-template",
    "tag": "latest",
    "environment": "production",
    "commit_sha": "abc123...",
    "timestamp": "2026-01-18T12:00:00Z"
  }
}
```

## Customization

To create a new service from this template:

1. Update `DISPLAY_MESSAGE` environment variables
2. Modify [templates/index.html](templates/index.html) for custom styling
3. Update service name in workflow file
4. Customize `main.go` for your specific needs
5. Adjust health check endpoints as needed

## Development

```bash
# Install dependencies
go mod download

# Run tests
go test -v ./...

# Run linter
golangci-lint run

# Build
go build -o service-template

# Build Docker image
docker build -t service-template .
```

## License

MIT License - feel free to use this template for your homelab projects!

## Contributing

This is a template repository. Fork it and make it your own!

