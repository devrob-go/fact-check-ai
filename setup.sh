#!/bin/bash

echo "🚀 Fact-Check Application Setup"
echo "================================"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

echo "✅ Docker and Docker Compose are installed"

# Create environment file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "⚠️  Please edit .env file with your configuration values:"
    echo "   - Google OAuth2 credentials"
    echo "   - OpenAI API key"
    echo "   - Database configuration"
    echo "   - JWT secret"
    echo ""
    echo "Press Enter when you've configured the .env file..."
    read
fi

# Check if .env file exists and has required values
if [ -f .env ]; then
    echo "🔍 Checking environment configuration..."
    
    # Check for required environment variables
    if grep -q "your-google-client-id" .env; then
        echo "⚠️  Please configure GOOGLE_CLIENT_ID in .env file"
    fi
    
    if grep -q "your-google-client-secret" .env; then
        echo "⚠️  Please configure GOOGLE_CLIENT_SECRET in .env file"
    fi
    
    if grep -q "your-openai-api-key" .env; then
        echo "⚠️  Please configure OPENAI_API_KEY in .env file"
    fi
    
    if grep -q "your-secret-key-change-in-production" .env; then
        echo "⚠️  Please configure JWT_SECRET in .env file"
    fi
fi

echo ""
echo "🏗️  Building Docker images..."
docker-compose build

echo ""
echo "🚀 Starting the application..."
docker-compose up -d

echo ""
echo "⏳ Waiting for services to start..."
sleep 10

# Check if services are running
if docker-compose ps | grep -q "Up"; then
    echo ""
    echo "✅ Application is running!"
    echo "🌐 Frontend: http://localhost:3000"
    echo "🔧 Backend: http://localhost:8080"
    echo "🗄️  Database: localhost:5432"
    echo ""
    echo "📋 Useful commands:"
    echo "   docker-compose logs -f          # View all logs"
    echo "   docker-compose logs -f backend  # View backend logs"
    echo "   docker-compose logs -f frontend # View frontend logs"
    echo "   docker-compose down             # Stop the application"
    echo "   docker-compose up -d            # Start the application"
    echo ""
    echo "🎉 Setup complete! Open http://localhost:3000 in your browser."
else
    echo "❌ Failed to start services. Check logs with: docker-compose logs"
    exit 1
fi
