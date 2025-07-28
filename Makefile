.PHONY: setup dev clean build test lint docker-up docker-down db-migrate db-seed

# Setup development environment
setup:
	@echo "Setting up CertiTrack development environment..."
	@cp .env.example .env
	@echo "Installing backend dependencies..."
	@cd backend && go mod tidy
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
	@echo "Starting Docker services..."
	@docker-compose up -d postgres redis mailhog
	@sleep 10
	@echo "Running database migrations..."
	@$(MAKE) db-migrate
	@echo "Seeding database..."
	@$(MAKE) db-seed
	@echo "Setup complete! Run 'make dev' to start development servers."

# Start development servers
dev:
	@echo "Starting development environment..."
	@docker-compose up -d postgres redis mailhog
	@echo "Starting backend server..."
	@cd backend && go run cmd/server/main.go &
	@echo "Starting frontend server..."
	@cd frontend && npm run dev

# Clean up
clean:
	@echo "Cleaning up..."
	@docker-compose down -v
	@cd backend && go clean
	@cd frontend && rm -rf .next node_modules

# Build applications
build-backend:
	@echo "Building backend..."
	@cd backend && go build -o bin/server cmd/server/main.go

build-frontend:
	@echo "Building frontend..."
	@cd frontend && npm run build

build: build-backend build-frontend

# Run tests
test-backend:
	@echo "Running backend tests..."
	@cd backend && go test -v -race -cover ./...

test-frontend:
	@echo "Running frontend tests..."
	@cd frontend && npm test

test: test-backend test-frontend

# Lint code
lint-backend:
	@echo "Linting backend code..."
	@cd backend && golangci-lint run

lint-frontend:
	@echo "Linting frontend code..."
	@cd frontend && npm run lint

lint: lint-backend lint-frontend

# Docker commands
docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

docker-build:
	@docker-compose build

# Database commands
db-migrate:
	@echo "Running database migrations..."
	@cd backend && go run cmd/migrate/main.go up

db-seed:
	@echo "Seeding database..."
	@cd backend && go run cmd/seed/main.go

db-reset:
	@echo "Resetting database..."
	@cd backend && go run cmd/migrate/main.go down
	@$(MAKE) db-migrate
	@$(MAKE) db-seed

# Help
help:
	@echo "Available commands:"
	@echo "  setup       - Set up development environment"
	@echo "  dev         - Start development servers"
	@echo "  clean       - Clean up build artifacts and containers"
	@echo "  build       - Build both backend and frontend"
	@echo "  test        - Run all tests"
	@echo "  lint        - Lint all code"
	@echo "  docker-up   - Start Docker services"
	@echo "  docker-down - Stop Docker services"
	@echo "  db-migrate  - Run database migrations"
	@echo "  db-seed     - Seed database with test data"
	@echo "  db-reset    - Reset database"