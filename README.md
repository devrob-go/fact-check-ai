# Fact-Check Microservice

A comprehensive fact-checking application built with Go backend and React frontend, featuring Google OAuth2 authentication and AI-powered news verification.

## Features

- **Backend**: Go REST API with JWT authentication, rate limiting, and secure middleware
- **Frontend**: React dashboard with Google OAuth2 integration
- **Database**: PostgreSQL for news storage and user management
- **AI Integration**: OpenAI API for automated fact-checking
- **Security**: JWT tokens, CORS, input validation, and rate limiting
- **Docker**: Containerized deployment for both frontend and backend
- **CI/CD**: GitHub Actions workflow for automated testing and deployment

## Architecture

```
├── backend/          # Go REST API service
├── frontend/         # React frontend application
├── docker/           # Docker configuration files
├── k8s/             # Kubernetes deployment manifests
├── .github/         # GitHub Actions workflows
└── docs/            # Documentation and API specs
```

## Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 15+
- OpenAI API key

### Environment Setup
1. Copy `.env.example` to `.env` and configure your environment variables
2. Set up your Google OAuth2 credentials
3. Configure your OpenAI API key

### Development
```bash
# Quick Setup (Recommended)
./setup.sh

# Manual Setup
# Backend
cd backend
go mod download
go run cmd/server/main.go

# Frontend
cd frontend
npm install
npm run dev

# Docker (full stack)
docker-compose up -d
```

### Production
```bash
# Build and deploy
docker-compose -f docker-compose.prod.yml up -d

# Kubernetes deployment
kubectl apply -f k8s/
```

## API Endpoints

- `POST /auth/login` - Google OAuth2 login
- `GET /auth/callback` - OAuth2 callback handler
- `POST /auth/logout` - User logout
- `POST /news/submit` - Submit news for verification
- `GET /news/verify/:id` - Verify news using AI
- `GET /news/user/:id` - Get user's news submissions

## 🚀 Quick Start

1. **Clone the repository**
   ```bash
   git clone git@github.com:devrob-go/go-chat-ai.git
   cd fact-check
   ```

2. **Run the setup script**
   ```bash
   ./setup.sh
   ```

3. **Configure your environment**
   - Edit `.env` file with your credentials
   - Set up Google OAuth2 credentials
   - Configure OpenAI API key

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend: http://localhost:8080

## 🔧 Troubleshooting

If you encounter issues, check the [Troubleshooting Guide](TROUBLESHOOTING.md) for common solutions.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests and ensure they pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details
