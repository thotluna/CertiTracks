# CertiTrack - Estructura del Proyecto y Guías de Desarrollo

## Visión General

Este documento define la estructura del proyecto, estándares de codificación, flujos de trabajo de desarrollo y mejores prácticas para el sistema de gestión de certificaciones CertiTrack, con el fin de garantizar consistencia, mantenibilidad y escalabilidad.

## Estructura del Proyecto

```
certitrack/
├── .github/                          # GitHub workflows and templates
│   ├── workflows/
│   │   ├── ci.yml                    # Continuous Integration
│   │   ├── cd.yml                    # Continuous Deployment
│   │   └── security-scan.yml         # Security scanning
│   ├── ISSUE_TEMPLATE/
│   └── PULL_REQUEST_TEMPLATE.md
├── docs/                             # Project documentation
│   ├── architecture/                 # Architecture documents
│   ├── api/                         # API documentation
│   ├── deployment/                  # Deployment guides
│   └── user-guides/                 # User documentation
├── scripts/                          # Utility scripts
│   ├── setup/                       # Setup scripts
│   ├── deployment/                  # Deployment scripts
│   ├── database/                    # Database utilities
│   └── monitoring/                  # Monitoring scripts
├── frontend/                         # Next.js frontend application
│   ├── public/                      # Static assets
│   │   ├── images/
│   │   ├── icons/
│   │   └── favicon.ico
│   ├── src/                         # Source code
│   │   ├── components/              # Reusable React components
│   │   │   ├── ui/                  # Basic UI components
│   │   │   ├── forms/               # Form components
│   │   │   ├── layout/              # Layout components
│   │   │   └── features/            # Feature-specific components
│   │   ├── pages/                   # Next.js pages
│   │   │   ├── api/                 # API routes (if needed)
│   │   │   ├── auth/                # Authentication pages
│   │   │   ├── dashboard/           # Dashboard pages
│   │   │   ├── people/              # People management
│   │   │   ├── equipment/           # Equipment management
│   │   │   ├── certifications/      # Certification management
│   │   │   └── reports/             # Reports pages
│   │   ├── hooks/                   # Custom React hooks
│   │   ├── services/                # API services and utilities
│   │   ├── utils/                   # Utility functions
│   │   ├── types/                   # TypeScript type definitions
│   │   ├── contexts/                # React contexts
│   │   ├── styles/                  # Global styles and themes
│   │   └── constants/               # Application constants
│   ├── tests/                       # Test files
│   │   ├── __mocks__/               # Mock files
│   │   ├── components/              # Component tests
│   │   ├── pages/                   # Page tests
│   │   ├── utils/                   # Utility tests
│   │   └── e2e/                     # End-to-end tests
│   ├── .env.example                 # Environment variables template
│   ├── .eslintrc.json              # ESLint configuration
│   ├── .prettierrc                 # Prettier configuration
│   ├── jest.config.js              # Jest configuration
│   ├── next.config.js              # Next.js configuration
│   ├── package.json                # Dependencies and scripts
│   ├── tailwind.config.js          # Tailwind CSS configuration
│   ├── tsconfig.json               # TypeScript configuration
│   └── Dockerfile                  # Docker configuration
├── backend/                          # Go backend application
│   ├── cmd/                         # Application entry points
│   │   ├── server/                  # Main server application
│   │   │   └── main.go
│   │   ├── worker/                  # Background workers
│   │   │   └── main.go
│   │   ├── migrate/                 # Database migration tool
│   │   │   └── main.go
│   │   └── seed/                    # Database seeding tool
│   │       └── main.go
│   ├── internal/                    # Private application code
│   │   ├── config/                  # Configuration management
│   │   │   ├── config.go
│   │   │   └── database.go
│   │   ├── handlers/                # HTTP handlers
│   │   │   ├── auth.go
│   │   │   ├── users.go
│   │   │   ├── people.go
│   │   │   ├── equipment.go
│   │   │   ├── certifications.go
│   │   │   ├── notifications.go
│   │   │   └── reports.go
│   │   ├── services/                # Business logic layer
│   │   │   ├── auth_service.go
│   │   │   ├── user_service.go
│   │   │   ├── person_service.go
│   │   │   ├── equipment_service.go
│   │   │   ├── certification_service.go
│   │   │   ├── notification_service.go
│   │   │   └── report_service.go
│   │   ├── repositories/            # Data access layer
│   │   │   ├── interfaces.go
│   │   │   ├── user_repository.go
│   │   │   ├── person_repository.go
│   │   │   ├── equipment_repository.go
│   │   │   ├── certification_repository.go
│   │   │   └── notification_repository.go
│   │   ├── models/                  # Data models
│   │   │   ├── user.go
│   │   │   ├── person.go
│   │   │   ├── equipment.go
│   │   │   ├── certification.go
│   │   │   ├── notification.go
│   │   │   └── audit.go
│   │   ├── middleware/              # HTTP middleware
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   ├── logging.go
│   │   │   ├── rate_limit.go
│   │   │   └── validation.go
│   │   ├── workers/                 # Background workers
│   │   │   ├── notification_worker.go
│   │   │   └── cleanup_worker.go
│   │   └── utils/                   # Utility functions
│   │       ├── crypto.go
│   │       ├── validation.go
│   │       └── response.go
│   ├── pkg/                         # Public packages
│   │   ├── database/                # Database utilities
│   │   │   ├── connection.go
│   │   │   └── migration.go
│   │   ├── email/                   # Email utilities
│   │   │   └── smtp.go
│   │   ├── storage/                 # File storage utilities
│   │   │   └── local.go
│   │   └── logger/                  # Logging utilities
│   │       └── logger.go
│   ├── migrations/                  # Database migrations
│   │   ├── 001_initial_schema.up.sql
│   │   ├── 001_initial_schema.down.sql
│   │   ├── 002_add_indexes.up.sql
│   │   └── 002_add_indexes.down.sql
│   ├── testdata/                    # Test data and fixtures
│   │   ├── fixtures/
│   │   └── seeds/
│   ├── tests/                       # Test files
│   │   ├── integration/             # Integration tests
│   │   ├── unit/                    # Unit tests
│   │   └── mocks/                   # Mock implementations
│   ├── .env.example                 # Environment variables template
│   ├── .golangci.yml               # Go linter configuration
│   ├── go.mod                      # Go module definition
│   ├── go.sum                      # Go module checksums
│   ├── Dockerfile                  # Docker configuration
│   └── Makefile                    # Build automation
├── database/                         # Database related files
│   ├── init/                        # Database initialization
│   │   └── 01-init.sql
│   ├── seeds/                       # Seed data
│   │   ├── users.sql
│   │   ├── certification_types.sql
│   │   └── notification_rules.sql
│   └── backups/                     # Database backups
├── nginx/                           # Nginx configuration
│   ├── nginx.conf                  # Main configuration
│   ├── ssl/                        # SSL certificates
│   └── conf.d/                     # Additional configurations
├── monitoring/                      # Monitoring and observability
│   ├── prometheus/                 # Prometheus configuration
│   ├── grafana/                    # Grafana dashboards
│   └── alerts/                     # Alert configurations
├── .env.example                     # Global environment template
├── .gitignore                      # Git ignore rules
├── .dockerignore                   # Docker ignore rules
├── docker-compose.yml              # Development Docker setup
├── docker-compose.prod.yml         # Production Docker setup
├── Makefile                        # Project automation
├── README.md                       # Project documentation
└── LICENSE                         # License file
```

## Estándares de Código

### Estándares del Backend en Go

#### Organización del Código

```go
// Package declaration and imports
package services

import (
    "context"
    "fmt"
    "time"

    // Standard library imports first
    "github.com/google/uuid"
    "gorm.io/gorm"

    // Third-party imports
    "github.com/gin-gonic/gin"

    // Local imports last
    "certitrack/internal/models"
    "certitrack/internal/repositories"
)
```

#### Naming Conventions

```go
// Constants: UPPER_SNAKE_CASE
const (
    DEFAULT_PAGE_SIZE = 20
    MAX_FILE_SIZE     = 10 * 1024 * 1024
)

// Variables and functions: camelCase
var defaultConfig = Config{}

func getUserByEmail(email string) (*User, error) {
    // Implementation
}

// Types: PascalCase
type UserService struct {
    repo repositories.UserRepository
    log  *logrus.Logger
}

// Interfaces: PascalCase with descriptive names
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

#### Manejo de Errores

```go
// Define custom error types
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidEmail     = errors.New("invalid email format")
    ErrDuplicateEmail   = errors.New("email already exists")
)

// Wrap errors with context
func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
    if err := s.validateUser(user); err != nil {
        return fmt.Errorf("user validation failed: %w", err)
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return nil
}

// Handle errors appropriately in handlers
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request format"})
        return
    }

    user, err := h.service.CreateUser(c.Request.Context(), &req.User)
    if err != nil {
        switch {
        case errors.Is(err, ErrDuplicateEmail):
            c.JSON(409, gin.H{"error": "Email already exists"})
        case errors.Is(err, ErrInvalidEmail):
            c.JSON(400, gin.H{"error": "Invalid email format"})
        default:
            h.logger.WithError(err).Error("Failed to create user")
            c.JSON(500, gin.H{"error": "Internal server error"})
        }
        return
    }

    c.JSON(201, gin.H{"data": user})
}
```

#### Estándares de Pruebas

```go
// Test file naming: *_test.go
// Test function naming: TestFunctionName_Scenario

func TestUserService_CreateUser_Success(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.UserRepository)
    service := services.NewUserService(mockRepo, nil)
    
    user := &models.User{
        Email:     "test@example.com",
        FirstName: "John",
        LastName:  "Doe",
    }

    mockRepo.On("Create", mock.Anything, user).Return(nil)

    // Act
    err := service.CreateUser(context.Background(), user)

    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

// Table-driven tests for multiple scenarios
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name          string
        input         *models.User
        setupMocks    func(*mocks.UserRepository)
        expectedError string
    }{
        {
            name: "successful creation",
            input: &models.User{Email: "test@example.com"},
            setupMocks: func(m *mocks.UserRepository) {
                m.On("Create", mock.Anything, mock.Anything).Return(nil)
            },
            expectedError: "",
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Estándares de Frontend (TypeScript/React)

#### Organización de Archivos

```typescript
// Component file structure
// components/UserForm/index.ts - Export file
// components/UserForm/UserForm.tsx - Main component
// components/UserForm/UserForm.test.tsx - Tests
// components/UserForm/UserForm.stories.tsx - Storybook stories
// components/UserForm/types.ts - Type definitions
```

#### Naming Conventions

```typescript
// Components: PascalCase
export const UserForm: React.FC<UserFormProps> = ({ onSubmit }) => {
  // Implementation
};

// Hooks: camelCase starting with 'use'
export const useUserData = (userId: string) => {
  // Implementation
};

// Types and Interfaces: PascalCase
interface UserFormProps {
  onSubmit: (user: User) => void;
  initialData?: Partial<User>;
}

type UserStatus = 'active' | 'inactive' | 'pending';

// Constants: UPPER_SNAKE_CASE
const API_ENDPOINTS = {
  USERS: '/api/v1/users',
  CERTIFICATIONS: '/api/v1/certifications',
} as const;
```

#### Estructura de Componentes

```typescript
// Component template
import React, { useState, useEffect } from 'react';
import { User } from '@/types/user';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';

interface UserFormProps {
  onSubmit: (user: User) => void;
  initialData?: Partial<User>;
  isLoading?: boolean;
}

export const UserForm: React.FC<UserFormProps> = ({
  onSubmit,
  initialData,
  isLoading = false,
}) => {
  // State declarations
  const [formData, setFormData] = useState<Partial<User>>(initialData || {});
  const [errors, setErrors] = useState<Record<string, string>>({});

  // Effects
  useEffect(() => {
    if (initialData) {
      setFormData(initialData);
    }
  }, [initialData]);

  // Event handlers
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    const validationErrors = validateForm(formData);
    if (Object.keys(validationErrors).length > 0) {
      setErrors(validationErrors);
      return;
    }

    onSubmit(formData as User);
  };

  const handleInputChange = (field: keyof User, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    
    // Clear error when user starts typing
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  // Helper functions
  const validateForm = (data: Partial<User>): Record<string, string> => {
    const errors: Record<string, string> = {};
    
    if (!data.email) {
      errors.email = 'Email is required';
    } else if (!isValidEmail(data.email)) {
      errors.email = 'Invalid email format';
    }

    return errors;
  };

  // Render
  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <Input
        label="Email"
        type="email"
        value={formData.email || ''}
        onChange={(value) => handleInputChange('email', value)}
        error={errors.email}
        required
      />
      
      <Button
        type="submit"
        disabled={isLoading}
        loading={isLoading}
      >
        {isLoading ? 'Saving...' : 'Save User'}
      </Button>
    </form>
  );
};
```

#### Custom Hooks

```typescript
// Custom hook example
import { useState, useEffect } from 'react';
import { User } from '@/types/user';
import { apiClient } from '@/services/api';

interface UseUserDataReturn {
  user: User | null;
  loading: boolean;
  error: string | null;
  refetch: () => void;
}

export const useUserData = (userId: string): UseUserDataReturn => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchUser = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await apiClient.get<User>(`/users/${userId}`);
      setUser(response.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch user');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (userId) {
      fetchUser();
    }
  }, [userId]);

  return {
    user,
    loading,
    error,
    refetch: fetchUser,
  };
};
```

## Flujo de Trabajo de Desarrollo

### Flujo de Trabajo con Git

#### Estrategia de Ramas

```bash
# Ramas principales
main          # Código listo para producción
develop       # Rama de integración para características

# Ramas de características
feature/user-management
feature/notification-system
feature/file-upload

# Ramas de lanzamiento
release/v1.0.0
release/v1.1.0

# Ramas de corrección
hotfix/security-patch
hotfix/critical-bug-fix
```

#### Convención de Mensajes de Commit

```bash
# Formato: <tipo>(<ámbito>): <descripción>

# Tipos:
feat:     # Nueva característica
fix:      # Corrección de error
docs:     # Cambios en la documentación
style:    # Cambios en el estilo del código (formato, etc.)
refactor: # Refactorización de código
test:     # Adición o actualización de pruebas
chore:    # Tareas de mantenimiento

# Ejemplos:
feat(auth): agregar mecanismo de actualización de token JWT
fix(notifications): resolver problema de renderizado de plantilla de correo
docs(api): actualizar documentación del punto final de autenticación
test(users): agregar pruebas unitarias para el servicio de usuarios
refactor(database): optimizar consultas de certificaciones
```

#### Proceso de Pull Request

```markdown
## Plantilla de Pull Request

### Descripción
Breve descripción de los cambios realizados.

### Tipo de Cambio
- [ ] Corrección de error (cambio no problemático que soluciona un problema)
- [ ] Nueva característica (cambio no problemático que agrega funcionalidad)
- [ ] Cambio problemático (solución o característica que haría que la funcionalidad existente no funcione como se espera)
- [ ] Actualización de documentación

### Pruebas
- [ ] Pruebas unitarias aprobadas
- [ ] Pruebas de integración aprobadas
- [ ] Pruebas de extremo a extremo aprobadas
- [ ] Pruebas manuales completadas

### Seguridad
- [ ] Sin datos sensibles expuestos
- [ ] Validación de entrada implementada
- [ ] Verificaciones de autorización en su lugar

### Lista de Verificación
- [ ] El código sigue las guías de estilo del proyecto
- [ ] Autorevisión completada
- [ ] Documentación actualizada
- [ ] Sin console.log o declaraciones de depuración
```

### Directrices para la Revisión de Código

#### Lista de Verificación para Revisión

**Funcionalidad**
- [ ] El código resuelve el problema previsto
- [ ] Se manejan los casos extremos
- [ ] El manejo de errores es apropiado
- [ ] Se han considerado aspectos de rendimiento

**Seguridad**
- [ ] Validación de entrada implementada
- [ ] Sin vulnerabilidades de inyección SQL
- [ ] Verificaciones de autenticación/autorización
- [ ] Manejo adecuado de datos sensibles

**Calidad del Código**
- [ ] El código es legible y está bien documentado
- [ ] Las funciones están enfocadas y tienen un único propósito
- [ ] Se siguen las convenciones de nomenclatura
- [ ] No hay duplicación de código

**Pruebas**
- [ ] Cobertura de pruebas adecuada
- [ ] Las pruebas son significativas y completas
- [ ] Uso apropiado de mocks

### Entorno de Desarrollo

#### Herramientas Requeridas

```bash
# Desarrollo de Backend
go version go1.21+
golangci-lint
migrate
air (para recarga en caliente)

# Desarrollo de Frontend
node v18+
npm v9+
eslint
prettier
typescript

# Base de Datos
postgresql 15+
redis 7+

# Contenedores
docker
docker-compose

# Control de Versiones
git
```

#### Configuración del IDE

**Extensiones de VS Code**
- Go (golang.go)
- TypeScript y JavaScript (ms-vscode.vscode-typescript-next)
- ESLint (dbaeumer.vscode-eslint)
- Prettier (esbenp.prettier-vscode)
- Tailwind CSS IntelliSense (bradlc.vscode-tailwindcss)
- GitLens (eamodio.gitlens)

**Configuración**
```json
{
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true,
    "source.fixAll.eslint": true
  },
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "typescript.preferences.importModuleSpecifier": "relative"
}
```

## Construcción y Despliegue

### Comandos de Makefile

```makefile
# Comandos de Backend
.PHONY: build-backend test-backend lint-backend

build-backend:
	cd backend && go build -o bin/server cmd/server/main.go

test-backend:
	cd backend && go test -v -race -cover ./...

lint-backend:
	cd backend && golangci-lint run

# Comandos de Frontend
.PHONY: build-frontend test-frontend lint-frontend

build-frontend:
	cd frontend && npm run build

test-frontend:
	cd frontend && npm test

lint-frontend:
	cd frontend && npm run lint

# Comandos de Docker
.PHONY: docker-build docker-up docker-down

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Comandos de Base de Datos
.PHONY: db-migrate db-seed db-reset

db-migrate:
	cd backend && migrate -path migrations -database "$(DB_URL)" up

db-seed:
	cd backend && go run cmd/seed/main.go

db-reset:
	cd backend && migrate -path migrations -database "$(DB_URL)" drop -f
	$(MAKE) db-migrate
	$(MAKE) db-seed

# Comandos de Desarrollo
.PHONY: dev setup clean

dev:
	docker-compose up -d postgres redis
	cd backend && air &
	cd frontend && npm run dev

setup:
	cd backend && go mod download
	cd frontend && npm install
	$(MAKE) db-migrate
	$(MAKE) db-seed

clean:
	docker-compose down -v
	cd backend && go clean
	cd frontend && rm -rf .next node_modules
```

### Pipeline de CI/CD

```yaml
# .github/workflows/ci.yml
name: CI Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: Install dependencies
        run: cd backend && go mod download
      
      - name: Run linter
        run: cd backend && golangci-lint run
      
      - name: Run tests
        run: cd backend && go test -v -race -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json
      
      - name: Install dependencies
        run: cd frontend && npm ci
      
      - name: Run linter
        run: cd frontend && npm run lint
      
      - name: Run tests
        run: cd frontend && npm run test:coverage
      
      - name: Build
        run: cd frontend && npm run build
```

Esta estructura de proyecto y guías de desarrollo proporcionan una base sólida para construir y mantener CertiTrack con consistencia, calidad y escalabilidad en mente.