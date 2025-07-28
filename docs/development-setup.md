# CertiTrack - Guía de Configuración del Entorno de Desarrollo

## Visión General

Esta guía proporciona instrucciones paso a paso para configurar un entorno de desarrollo local para CertiTrack, incluyendo todas las herramientas, dependencias y configuraciones necesarias.

## Requisitos Previos

### Requisitos del Sistema

- **Sistema Operativo**: macOS, Linux o Windows con WSL2
- **RAM**: Mínimo 8GB, recomendado 16GB
- **Almacenamiento**: Al menos 10GB de espacio libre
- **Red**: Conexión a Internet estable para descargar dependencias

### Software Requerido

1. **Docker & Docker Compose**
   - Docker Desktop 4.0+ (incluye Docker Compose)
   - Alternativa: Docker Engine + plugin de Docker Compose

2. **Git**
   - Versión 2.30 o superior

3. **Node.js & npm**
   - Node.js 18.x LTS
   - npm 9.x

4. **Go**
   - Go 1.21 o superior

5. **Editor de Código**
   - VS Code (recomendado) con extensiones
   - Alternativas: GoLand, WebStorm o cualquier editor preferido

## Guía de Instalación

### 1. Instalar Docker

#### macOS
```bash
# Install Docker Desktop
brew install --cask docker

# Start Docker Desktop
open /Applications/Docker.app
```

#### Linux (Ubuntu/Debian)
```bash
# Update package index
sudo apt-get update

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt-get install docker-compose-plugin

# Restart to apply group changes
sudo reboot
```

#### Windows
1. Install Docker Desktop from https://www.docker.com/products/docker-desktop
2. Enable WSL2 integration
3. Restart system

### 2. Instalar Herramientas de Desarrollo

#### Node.js
```bash
# Using Node Version Manager (recommended)
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
source ~/.bashrc
nvm install 18
nvm use 18

# Verify installation
node --version
npm --version
```

#### Go
```bash
# macOS
brew install go

# Linux
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

### 3. Clonar el Repositorio

```bash
# Clone the repository
git clone https://github.com/company/certitrack.git
cd certitrack

# Create development branch
git checkout -b develop
```

## Estructura del Proyecto

```
certitrack/
├── frontend/                 # Next.js frontend application
│   ├── src/
│   │   ├── components/      # Reusable React components
│   │   ├── pages/          # Next.js pages
│   │   ├── hooks/          # Custom React hooks
│   │   ├── services/       # API services
│   │   ├── utils/          # Utility functions
│   │   └── styles/         # CSS/SCSS styles
│   ├── public/             # Static assets
│   ├── package.json
│   ├── next.config.js
│   └── Dockerfile
├── backend/                  # Go API server
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   │   ├── handlers/       # HTTP handlers
│   │   ├── services/       # Business logic
│   │   ├── repositories/   # Data access layer
│   │   ├── models/         # Data models
│   │   ├── middleware/     # HTTP middleware
│   │   └── config/         # Configuration
│   ├── pkg/                # Public packages
│   ├── migrations/         # Database migrations
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile
├── database/                # Database related files
│   ├── migrations/         # SQL migration files
│   ├── seeds/             # Test data
│   └── init/              # Initialization scripts
├── nginx/                   # Nginx configuration
│   ├── nginx.conf
│   └── ssl/
├── scripts/                 # Utility scripts
│   ├── setup.sh
│   ├── migrate.sh
│   └── seed.sh
├── docs/                    # Documentation
├── docker-compose.yml       # Development Docker setup
├── docker-compose.prod.yml  # Production Docker setup
├── .env.example            # Environment variables template
├── .gitignore
└── README.md
```

## Configuración del Entorno de Desarrollo

### 1. Configuración del Entorno

```bash
# Copy environment template
cp .env.example .env.development

# Edit environment variables
nano .env.development
```

#### .env.development
```env
# Application
APP_ENV=development
APP_URL=http://localhost:3000
API_URL=http://localhost:8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=certitrack_dev
DB_USER=certitrack_user
DB_PASSWORD=dev_password

# Redis
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=dev_redis_password

# JWT
JWT_SECRET=development-jwt-secret-key-minimum-32-characters
JWT_ACCESS_TOKEN_EXPIRY=15m
JWT_REFRESH_TOKEN_EXPIRY=168h

# SMTP (for development - use Mailhog)
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM=CertiTrack Dev <dev@certitrack.local>

# File Storage
STORAGE_ROOT=./storage
MAX_FILE_SIZE_MB=10

# Logging
LOG_LEVEL=debug
ENABLE_METRICS=true

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_NAME=CertiTrack
```

### 2. Configuración de Docker para Desarrollo

#### docker-compose.yml (Development)
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: certitrack-postgres-dev
    environment:
      POSTGRES_DB: certitrack_dev
      POSTGRES_USER: certitrack_user
      POSTGRES_PASSWORD: dev_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_dev_data:/var/lib/postgresql/data
      - ./database/init:/docker-entrypoint-initdb.d
    networks:
      - certitrack-dev

  redis:
    image: redis:7-alpine
    container_name: certitrack-redis-dev
    command: redis-server --requirepass dev_redis_password
    ports:
      - "6379:6379"
    volumes:
      - redis_dev_data:/data
    networks:
      - certitrack-dev

  mailhog:
    image: mailhog/mailhog:latest
    container_name: certitrack-mailhog-dev
    ports:
      - "1025:1025"  # SMTP
      - "8025:8025"  # Web UI
    networks:
      - certitrack-dev

volumes:
  postgres_dev_data:
  redis_dev_data:

networks:
  certitrack-dev:
    driver: bridge
```

### 3. Iniciar Servicios de Desarrollo

```bash
# Start database and supporting services
docker-compose up -d

# Verify services are running
docker-compose ps

# Check logs if needed
docker-compose logs postgres
```

### 4. Configuración del Backend

```bash
cd backend

# Install Go dependencies
go mod download

# Install development tools
go install github.com/cosmtrek/air@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run database migrations
migrate -path migrations -database "postgres://certitrack_user:dev_password@localhost:5432/certitrack_dev?sslmode=disable" up

# Seed database with test data
go run cmd/seed/main.go

# Start development server with hot reload
air
```

#### Air Configuration (.air.toml)
```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false
```

### 5. Configuración del Frontend

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Alternative: Start with specific port
npm run dev -- --port 3000
```

#### package.json Scripts
```json
{
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "lint:fix": "next lint --fix",
    "type-check": "tsc --noEmit",
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage"
  }
}
```

## Configuración de VS Code

### Extensiones Recomendadas

Create `.vscode/extensions.json`:
```json
{
  "recommendations": [
    "golang.go",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode",
    "ms-vscode.vscode-typescript-next",
    "ms-vscode.vscode-json",
    "redhat.vscode-yaml",
    "ms-vscode-remote.remote-containers",
    "github.copilot"
  ]
}
```

### Configuración del Espacio de Trabajo

Create `.vscode/settings.json`:
```json
{
  "go.toolsManagement.checkForUpdates": "local",
  "go.useLanguageServer": true,
  "go.gopath": "",
  "go.goroot": "",
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  },
  "typescript.preferences.importModuleSpecifier": "relative",
  "eslint.workingDirectories": ["frontend"],
  "prettier.configPath": "frontend/.prettierrc"
}
```

### Configuración de Lanzamiento

Create `.vscode/launch.json`:
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Backend",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/backend/cmd/server",
      "env": {
        "APP_ENV": "development"
      },
      "args": []
    },
    {
      "name": "Debug Test",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${workspaceFolder}/backend"
    }
  ]
}
```

## Flujo de Trabajo de Desarrollo

### 1. Rutina Diaria de Desarrollo

```bash
# Start development environment
./scripts/dev-start.sh

# Pull latest changes
git pull origin develop

# Update dependencies if needed
cd backend && go mod tidy
cd ../frontend && npm install

# Start coding!
```

### 2. Scripts de Desarrollo

Create `scripts/dev-start.sh`:
```bash
#!/bin/bash
set -e

echo "Starting CertiTrack development environment..."

# Start supporting services
docker-compose up -d postgres redis mailhog

# Wait for services to be ready
echo "Waiting for services to start..."
sleep 10

# Check if database needs migration
cd backend
if ! migrate -path migrations -database "postgres://certitrack_user:dev_password@localhost:5432/certitrack_dev?sslmode=disable" version; then
    echo "Running database migrations..."
    migrate -path migrations -database "postgres://certitrack_user:dev_password@localhost:5432/certitrack_dev?sslmode=disable" up
fi

echo "Development environment ready!"
echo "Backend: http://localhost:8080"
echo "Frontend: http://localhost:3000"
echo "Mailhog: http://localhost:8025"
echo "Database: localhost:5432"
```

Create `scripts/dev-stop.sh`:
```bash
#!/bin/bash
echo "Stopping CertiTrack development environment..."
docker-compose down
echo "Development environment stopped."
```

### 3. Gestión de la Base de Datos

```bash
# Create new migration
migrate create -ext sql -dir backend/migrations -seq add_new_table

# Run migrations
migrate -path backend/migrations -database "postgres://certitrack_user:dev_password@localhost:5432/certitrack_dev?sslmode=disable" up

# Rollback migration
migrate -path backend/migrations -database "postgres://certitrack_user:dev_password@localhost:5432/certitrack_dev?sslmode=disable" down 1

# Reset database
migrate -path backend/migrations -database "postgres://certitrack_user:dev_password@localhost:5432/certitrack_dev?sslmode=disable" drop -f
migrate -path backend/migrations -database "postgres://certitrack_user:dev_password@localhost:5432/certitrack_dev?sslmode=disable" up
```

## Configuración de Pruebas

### Pruebas del Backend

```bash
cd backend

# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestUserService ./internal/services
```

### Pruebas del Frontend

```bash
cd frontend

# Run unit tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage
npm run test:coverage

# Run E2E tests (if configured)
npm run test:e2e
```

## Depuración

### Depuración del Backend

1. **Using VS Code Debugger**
   - Set breakpoints in code
   - Press F5 to start debugging
   - Use Debug Console for evaluation

2. **Using Delve (command line)**
   ```bash
   # Install delve
   go install github.com/go-delve/delve/cmd/dlv@latest
   
   # Debug application
   dlv debug ./cmd/server
   ```

### Depuración del Frontend

1. **Browser DevTools**
   - Use Chrome/Firefox DevTools
   - React Developer Tools extension

2. **VS Code Debugging**
   - Install "Debugger for Chrome" extension
   - Configure launch.json for browser debugging

## Problemas Comunes y Soluciones

### 1. Puerto en Uso
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>
```

### 2. Problemas de Permisos en Docker (Linux)
```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Restart session or reboot
```

### 3. Problemas de Conexión a la Base de Datos
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Reset database
docker-compose down -v
docker-compose up -d postgres
```

### 4. Problemas con Módulos de Go
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
```

### 5. Problemas con Node.js
```bash
# Clear npm cache
npm cache clean --force

# Delete node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

## Optimización de Rendimiento

### Consejos de Rendimiento para Desarrollo

1. **Use Air for Go hot reload**
2. **Enable Next.js Fast Refresh**
3. **Use Docker BuildKit for faster builds**
4. **Configure IDE for optimal performance**

### Optimización de Docker

```bash
# Enable BuildKit
export DOCKER_BUILDKIT=1

# Use multi-stage builds
# Optimize layer caching
# Use .dockerignore files
```

## Consideraciones de Seguridad

### Seguridad en Desarrollo

1. **Never commit secrets to git**
2. **Use different credentials for development**
3. **Keep development dependencies updated**
4. **Use HTTPS in development when possible**

### Aislamiento del Entorno

```bash
# Use separate databases for different features
# Use Docker networks for service isolation
# Implement proper CORS settings
```

Esta configuración de desarrollo proporciona una base completa para construir CertiTrack de manera eficiente, manteniendo la consistencia en todo el equipo de desarrollo.