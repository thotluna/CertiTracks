.PHONY: setup dev stop clean build test lint docker-up docker-down db-migrate db-seed

# Setup development environment
setup:
	@echo "Setting up CertiTrack development environment..."
	@cp .env.example .env
	@echo "Installing backend dependencies..."
	@cd backend && go mod tidy
	@echo "Installing frontend dependencies..."
	@cd frontend && pnpm install
	@echo "Starting Docker services..."
	@docker-compose up -d postgres redis mailhog
	@sleep 10
	# @echo "Running database migrations..."
	# @$(MAKE) db-migrate
	# @echo "Seeding database..."
	# @$(MAKE) db-seed
	@echo "Setup complete! Run 'make dev' to start development servers."

# Find process using port 8080
PORT := 8080

# Start development servers
dev: stop
	@echo "Starting development environment..."
	@docker-compose up -d postgres redis mailhog
	@echo "Starting backend server..."
	@cd backend && (go run cmd/server/main.go & echo $$! > .server.pid)
	@sleep 1
	@echo "Backend server started on port $(PORT)"
	# @echo "Starting frontend server..."
	# @cd frontend && pnpm run dev

# Stop development servers
stop:
	@echo "Stopping development environment..."
	@docker-compose down
	@if lsof -ti:$(PORT) >/dev/null 2>&1; then \
		echo "Stopping process on port $(PORT)..."; \
		lsof -ti:$(PORT) | xargs kill -9 2>/dev/null || true; \
	fi
	@if [ -f "backend/.server.pid" ]; then \
		rm -f backend/.server.pid backend/.server.child.pid; \
	fi

# Clean up
clean: stop
	@echo "Cleaning up..."
	@docker-compose down -v
	@cd backend && go clean
	@cd frontend && rm -rf .next node_modules
	@rm -f backend/.server.pid backend/.server.child.pid 2>/dev/null || true

# Build applications
build-backend: wire-gen
	@echo "Building backend..."
	@cd backend && go build -o bin/server cmd/server/main.go

build-frontend:
	@echo "Building frontend..."
	@cd frontend && pnpm run build

build: build-backend build-frontend

# Clean test cache
test-clean:
	@echo "Cleaning test cache..."
	@cd backend && go clean -testcache

test-backend: test-clean test-bu test-bi

test-bu: test-clean
	@echo "Running backend unit tests..."
	@cd backend && APP_ENV=test go test -short -v ./internal/...

test-bi: test-clean
	@echo "Running backend integration tests..."
	@cd backend && APP_ENV=test go test -run Integration -v ./integration/...

test-frontend:
	@echo "Running frontend tests..."
	@cd frontend && pnpm test

test: test-backend test-frontend

# Lint code
lint-backend:
	@echo "Linting backend code..."
	@cd backend && golangci-lint run

lint-frontend:
	@echo "Linting frontend code..."
	@cd frontend && pnpm run lint

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

wire-gen:
	@echo "Generating Wire dependencies..."
	@cd backend && wire ./internal/di

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