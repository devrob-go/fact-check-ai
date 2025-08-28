.PHONY: help install-backend install-frontend build run-docker stop-docker clean test-backend test-frontend deploy-k8s

help: ## Show this help message
	@echo "Fact-Check Application - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install-backend: ## Install Go backend dependencies
	@echo "Installing Go backend dependencies..."
	cd backend && go mod download
	cd backend && go mod tidy

install-frontend: ## Install React frontend dependencies
	@echo "Installing React frontend dependencies..."
	cd frontend && npm install

install: install-backend install-frontend ## Install all dependencies

build-backend: ## Build Go backend
	@echo "Building Go backend..."
	cd backend && go build -o bin/server ./cmd/server

build-frontend: ## Build React frontend
	@echo "Building React frontend..."
	cd frontend && npm run build

build: build-backend build-frontend ## Build both backend and frontend

run-backend: ## Run Go backend locally
	@echo "Running Go backend..."
	cd backend && go run ./cmd/server

run-frontend: ## Run React frontend locally
	@echo "Running React frontend..."
	cd frontend && npm start

run-docker: ## Run application with Docker Compose
	@echo "Starting application with Docker Compose..."
	docker-compose up -d

stop-docker: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	docker-compose down

run-docker-prod: ## Run production application with Docker Compose
	@echo "Starting production application with Docker Compose..."
	docker-compose -f docker-compose.prod.yml up -d

stop-docker-prod: ## Stop production Docker Compose services
	@echo "Stopping production Docker Compose services..."
	docker-compose -f docker-compose.prod.yml down

test-backend: ## Run Go backend tests
	@echo "Running Go backend tests..."
	cd backend && go test -v ./...

test-frontend: ## Run React frontend tests
	@echo "Running React frontend tests..."
	cd frontend && npm test -- --watchAll=false

test: test-backend test-frontend ## Run all tests

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin/
	rm -rf frontend/build/
	rm -rf frontend/node_modules/
	docker-compose down -v
	docker system prune -f

deploy-k8s: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/postgres-secret.yaml
	kubectl apply -f k8s/app-secret.yaml
	kubectl apply -f k8s/postgres-pvc.yaml
	kubectl apply -f k8s/postgres-deployment.yaml
	kubectl apply -f k8s/backend-deployment.yaml
	kubectl apply -f k8s/frontend-deployment.yaml

undeploy-k8s: ## Undeploy from Kubernetes
	@echo "Undeploying from Kubernetes..."
	kubectl delete -f k8s/ --ignore-not-found=true

logs-backend: ## View backend logs
	@echo "Viewing backend logs..."
	docker-compose logs -f backend

logs-frontend: ## View frontend logs
	@echo "Viewing frontend logs..."
	docker-compose logs -f frontend

logs-postgres: ## View PostgreSQL logs
	@echo "Viewing PostgreSQL logs..."
	docker-compose logs -f postgres

logs: ## View all logs
	@echo "Viewing all logs..."
	docker-compose logs -f

setup-dev: ## Setup development environment
	@echo "Setting up development environment..."
	cp env.example .env
	@echo "Please edit .env file with your configuration values"
	@echo "Then run: make install && make run-docker"

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	docker-compose build

docker-build-prod: ## Build production Docker images
	@echo "Building production Docker images..."
	docker-compose -f docker-compose.prod.yml build
