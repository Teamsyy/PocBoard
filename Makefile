.PHONY: help dev build clean stop logs test backend-deps frontend-deps

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development commands
dev: ## Start development environment with Docker Compose
	docker compose up --build --watch

dev-detached: ## Start development environment in detached mode
	docker compose up --build -d

stop: ## Stop all services
	docker compose down

clean: ## Stop services and remove volumes
	docker compose down -v --remove-orphans

restart: ## Restart all services
	docker compose restart

# Logs
logs: ## Show logs from all services
	docker compose logs -f

logs-backend: ## Show backend logs
	docker compose logs -f backend

logs-frontend: ## Show frontend logs
	docker compose logs -f frontend

logs-db: ## Show database logs
	docker compose logs -f postgres

# Database commands
db-reset: ## Reset database (remove volume and restart)
	docker compose down -v
	docker compose up postgres -d
	@echo "Waiting for database to be ready..."
	@sleep 10

db-shell: ## Connect to database shell
	docker compose exec postgres psql -U postgres -d junk_journal

# Backend commands
backend-shell: ## Open shell in backend container
	docker compose exec backend sh

backend-deps: ## Install backend dependencies
	cd backend && go mod tidy && go mod download

backend-test: ## Run backend tests
	cd backend && go test ./...

backend-build: ## Build backend binary
	cd backend && go build -o bin/main .

# Frontend commands
frontend-shell: ## Open shell in frontend container
	docker compose exec frontend sh

frontend-deps: ## Install frontend dependencies
	cd frontend && npm ci

frontend-test: ## Run frontend tests
	cd frontend && npm run test

frontend-lint: ## Run frontend linter
	cd frontend && npm run lint

frontend-type-check: ## Run TypeScript type checking
	cd frontend && npm run type-check

frontend-build: ## Build frontend for production
	cd frontend && npm run build

# Build commands
build: ## Build all services
	docker compose build

build-backend: ## Build backend service only
	docker compose build backend

build-frontend: ## Build frontend service only
	docker compose build frontend

# Production commands
prod-build: ## Build production images
	docker compose -f docker-compose.yml -f docker-compose.prod.yml build

prod-up: ## Start production environment
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Utility commands
ps: ## Show running containers
	docker compose ps

exec-backend: ## Execute command in backend container (usage: make exec-backend CMD="go version")
	docker compose exec backend $(CMD)

exec-frontend: ## Execute command in frontend container (usage: make exec-frontend CMD="npm --version")
	docker compose exec frontend $(CMD)

# Setup commands
setup: ## Initial setup - install dependencies and start services
	@echo "Setting up Junk Journal Board development environment..."
	@echo "Installing backend dependencies..."
	$(MAKE) backend-deps
	@echo "Installing frontend dependencies..."
	$(MAKE) frontend-deps
	@echo "Starting development environment..."
	$(MAKE) dev