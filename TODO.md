# Estado del Proyecto CertiTracks

## ✅ Funcionalidades Completadas

### Autenticación Básica
- [x] Registro de usuarios
- [x] Inicio de sesión con JWT
- [x] Refresh tokens
- [x] Cierre de sesión con revocación de tokens
- [x] Middleware de autenticación
- [x] Pruebas unitarias básicas

### Gestión de Perfil
- [x] Obtener perfil de usuario
- [x] Cambio de contraseña

## 🚧 En Progreso

### Recuperación de Contraseña
- [ ] Configuración inicial
  - [x] Modelo `PasswordResetToken`
  - [ ] Configuración de servicio de correo

- [ ] Solicitud de restablecimiento
  - [ ] Endpoint POST `/auth/forgot-password`
  - [ ] Pruebas unitarias
  - [ ] Integración con servicio de correo

- [ ] Validación de token
  - [ ] Endpoint GET `/auth/validate-reset-token/{token}`
  - [ ] Pruebas de validación

- [ ] Restablecimiento de contraseña
  - [ ] Endpoint POST `/auth/reset-password`
  - [ ] Pruebas de actualización
  - [ ] Notificaciones por correo

## 📅 Próximos Pasos (Priorizados)

### Prioridad Alta
1. **Recuperación de Contraseña**
   - Completar implementación de endpoints faltantes
   - Implementar servicio de correo electrónico
   - Añadir pruebas de integración

2. **Seguridad Mejorada**
   - [ ] Limpieza periódica de tokens expirados
   - [ ] Revocación masiva de tokens por usuario
   - [ ] Límite de intentos de inicio de sesión
   - [ ] Bloqueo temporal de cuentas

3. **Verificación de Correo Electrónico**
   - [ ] Modelo `EmailVerificationToken`
   - [ ] Endpoint para reenviar correo de verificación
   - [ ] Endpoint para validar token de verificación
   - [ ] Actualización de estado de verificación

### Prioridad Media
1. **Mejoras de Seguridad**
   - [ ] Registro de actividades sospechosas
   - [ ] Detección de patrones de acceso inusuales
   - [ ] Soporte para autenticación de dos factores (2FA)

2. **Gestión de Sesiones**
   - [ ] Listado de sesiones activas
   - [ ] Cierre de sesión remota
   - [ ] Historial de inicios de sesión

3. **API y Documentación**
   - [ ] Documentación Swagger/OpenAPI
   - [ ] Ejemplos de código para integración
   - [ ] Guía de implementación para frontend

### Prioridad Baja
1. **Mejoras de Usuario**
   - [ ] Perfil de usuario mejorado
   - [ ] Preferencias de notificación
   - [ ] Personalización de interfaz

2. **Monitoreo y Analíticas**
   - [ ] Métricas de uso de la API
   - [ ] Registro de eventos de seguridad
   - [ ] Panel de administración

## 🏗️ Implementación de Recursos: Persons y Equipments

### 1. Modelos de Datos (GORM)

#### User Model
```go
// internal/models/user.go
package models

type User struct {
    ID           uint   `gorm:"primaryKey" json:"id"`
    Email        string `gorm:"size:100;unique;not null" json:"email"`
    PasswordHash string `gorm:"not null" json:"-"`
    IsActive     bool   `gorm:"default:true" json:"is_active"`
    IsAdmin      bool   `gorm:"default:false" json:"is_admin"`
    LastLogin    *time.Time `json:"last_login,omitempty"`
    PersonID     *uint      `gorm:"index" json:"person_id,omitempty"`
    Person       *Person    `gorm:"foreignKey:PersonID" json:"person,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
```

#### Person Model
```go
// internal/models/person.go
package models

type Person struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    FirstName   string    `gorm:"size:100;not null" json:"first_name"`
    LastName    string    `gorm:"size:100;not null" json:"last_name"`
    Email       string    `gorm:"size:100;unique;not null" json:"email"`
    Phone       string    `gorm:"size:20" json:"phone"`
    Position    string    `gorm:"size:100" json:"position"`
    Department  string    `gorm:"size:100" json:"department"`
    IsActive    bool      `gorm:"default:true" json:"is_active"`
    HiredAt     time.Time `json:"hired_at"`
    User        *User     `gorm:"foreignKey:PersonID" json:"user,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

#### Equipment Model
```go
// internal/models/equipment.go
package models

type EquipmentType string

const (
    EquipmentTypeLaptop   EquipmentType = "laptop"
    EquipmentTypePhone    EquipmentType = "phone"
    EquipmentTypeTablet   EquipmentType = "tablet"
    EquipmentTypeOther    EquipmentType = "other"
)

type EquipmentStatus string

const (
    StatusAvailable   EquipmentStatus = "available"
    StatusAssigned    EquipmentStatus = "assigned"
    StatusMaintenance EquipmentStatus = "maintenance"
    StatusRetired     EquipmentStatus = "retired"
)

type Equipment struct {
    ID              uint            `gorm:"primaryKey" json:"id"`
    Name            string          `gorm:"size:100;not null" json:"name"`
    Type            EquipmentType   `gorm:"type:varchar(20);not null" json:"type"`
    Status          EquipmentStatus `gorm:"type:varchar(20);not null" json:"status"`
    SerialNumber    string          `gorm:"size:100;unique;not null" json:"serial_number"`
    Model           string          `gorm:"size:100" json:"model"`
    Manufacturer    string          `gorm:"size:100" json:"manufacturer"`
    PurchaseDate    time.Time       `json:"purchase_date"`
    PurchaseCost    float64         `json:"purchase_cost"`
    WarrantyExpires *time.Time      `json:"warranty_expires,omitempty"`
    AssignedTo      *uint           `json:"assigned_to,omitempty"`
    AssignedPerson  *Person         `gorm:"foreignKey:AssignedTo" json:"assigned_person,omitempty"`
    Notes           string          `gorm:"type:text" json:"notes"`
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
    DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
}
```

### 2. Relaciones entre Modelos

#### User Model (actualizado)
```go
// internal/models/user.go
package models

type User struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    Email        string     `gorm:"size:100;unique;not null" json:"email"`
    PasswordHash string     `gorm:"not null" json:"-"`
    IsActive     bool       `gorm:"default:true" json:"is_active"`
    IsAdmin      bool       `gorm:"default:false" json:"is_admin"`
    LastLogin    *time.Time `json:"last_login,omitempty"`
    
    // Relación con Person (opcional)
    PersonID     *uint      `gorm:"index" json:"person_id,omitempty"`
    Person       *Person    `gorm:"foreignKey:PersonID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"person,omitempty"`
    
    // Timestamps
    CreatedAt    time.Time       `json:"created_at"`
    UpdatedAt    time.Time       `json:"updated_at"`
    DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
}
```

#### Person Model (actualizado)
```go
// internal/models/person.go
package models

type Person struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    FirstName   string `gorm:"size:100;not null" json:"first_name"`
    LastName    string `gorm:"size:100;not null" json:"last_name"`
    Email       string `gorm:"size:100;unique;not null" json:"email"`
    Phone       string `gorm:"size:20" json:"phone"`
    Position    string `gorm:"size:100" json:"position"`
    Department  string `gorm:"size:100" json:"department"`
    IsActive    bool   `gorm:"default:true" json:"is_active"`
    HiredAt     time.Time `json:"hired_at"`
    
    // Relación con User (opcional)
    User        *User     `gorm:"foreignKey:PersonID" json:"user,omitempty"`
    
    // Relación con Equipment (una persona puede tener múltiples equipos)
    Equipment   []Equipment `gorm:"foreignKey:AssignedTo" json:"equipment,omitempty"`
    
    // Timestamps
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
```

#### Equipment Model (actualizado)
```go
// internal/models/equipment.go
package models

type EquipmentType string

const (
    EquipmentTypeLaptop   EquipmentType = "laptop"
    EquipmentTypePhone    EquipmentType = "phone"
    EquipmentTypeTablet   EquipmentType = "tablet"
    EquipmentTypeOther    EquipmentType = "other"
)

type EquipmentStatus string

const (
    StatusAvailable   EquipmentStatus = "available"
    StatusAssigned    EquipmentStatus = "assigned"
    StatusMaintenance EquipmentStatus = "maintenance"
    StatusRetired     EquipmentStatus = "retired"
)

type Equipment struct {
    ID              uint            `gorm:"primaryKey" json:"id"`
    Name            string          `gorm:"size:100;not null" json:"name"`
    Type            EquipmentType   `gorm:"type:varchar(20);not null" json:"type"`
    Status          EquipmentStatus `gorm:"type:varchar(20);not null;index" json:"status"`
    SerialNumber    string          `gorm:"size:100;unique;not null;index" json:"serial_number"`
    Model           string          `gorm:"size:100" json:"model"`
    Manufacturer    string          `gorm:"size:100" json:"manufacturer"`
    PurchaseDate    time.Time       `json:"purchase_date"`
    PurchaseCost    float64         `json:"purchase_cost"`
    WarrantyExpires *time.Time      `json:"warranty_expires,omitempty"`
    
    // Relación con Person (opcional)
    AssignedTo      *uint           `gorm:"index" json:"assigned_to,omitempty"`
    AssignedPerson  *Person         `gorm:"foreignKey:AssignedTo;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"assigned_person,omitempty"`
    
    Notes           string          `gorm:"type:text" json:"notes"`
    
    // Timestamps
    CreatedAt       time.Time       `json:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at"`
    DeletedAt       gorm.DeletedAt  `gorm:"index" json:"-"`
}
```

### 3. Inicialización de la Base de Datos

```go
// internal/database/database.go
package database

import (
    "gorm.io/gorm"
    "certitrack/internal/models"
)

func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Person{},
        &models.Equipment{},
    )
}
```

### 3. Repositorios

#### Person Repository
```go
// internal/repositories/person_repository.go
package repositories

type PersonRepository interface {
    Create(person *models.Person) error
    FindByID(id uint) (*models.Person, error)
    FindByEmail(email string) (*models.Person, error)
    Update(person *models.Person) error
    Delete(id uint) error
    List(filters map[string]interface{}) ([]models.Person, error)
}

type personRepository struct {
    db *gorm.DB
}

func NewPersonRepository(db *gorm.DB) PersonRepository {
    return &personRepository{db: db}
}

// Implementación de los métodos...
```

#### Equipment Repository
```go
// internal/repositories/equipment_repository.go
package repositories

type EquipmentRepository interface {
    Create(equipment *models.Equipment) error
    FindByID(id uint) (*models.Equipment, error)
    FindBySerialNumber(serialNumber string) (*models.Equipment, error)
    Update(equipment *models.Equipment) error
    Delete(id uint) error
    List(filters map[string]interface{}) ([]models.Equipment, error)
    GetAssignedEquipment(personID uint) ([]models.Equipment, error)
}

type equipmentRepository struct {
    db *gorm.DB
}

func NewEquipmentRepository(db *gorm.DB) EquipmentRepository {
    return &equipmentRepository{db: db}
}

// Implementación de los métodos...
```

### 4. Servicios

#### Person Service
```go
// internal/services/person_service.go
package services

type PersonService interface {
    RegisterPerson(person *models.Person) error
    GetPerson(id uint) (*models.Person, error)
    UpdatePerson(person *models.Person) error
    DeletePerson(id uint) error
    ListPeople(filters map[string]interface{}) ([]models.Person, error)
}

type personService struct {
    personRepo repositories.PersonRepository
}

func NewPersonService(repo repositories.PersonRepository) PersonService {
    return &personService{personRepo: repo}
}

// Implementación de los métodos...
```

#### Equipment Service
```go
// internal/services/equipment_service.go
package services

type EquipmentService interface {
    RegisterEquipment(equipment *models.Equipment) error
    GetEquipment(id uint) (*models.Equipment, error)
    UpdateEquipment(equipment *models.Equipment) error
    DeleteEquipment(id uint) error
    ListEquipment(filters map[string]interface{}) ([]models.Equipment, error)
    AssignEquipment(equipmentID, personID uint) error
    UnassignEquipment(equipmentID uint) error
    GetPersonEquipment(personID uint) ([]models.Equipment, error)
}

type equipmentService struct {
    equipmentRepo repositories.EquipmentRepository
    personRepo    repositories.PersonRepository
}

func NewEquipmentService(
    equipmentRepo repositories.EquipmentRepository,
    personRepo repositories.PersonRepository,
) EquipmentService {
    return &equipmentService{
        equipmentRepo: equipmentRepo,
        personRepo:    personRepo,
    }
}

// Implementación de los métodos...
```

### 5. Controladores (Handlers)

#### Person Handlers
```go
// internal/handlers/person_handlers.go
package handlers

type PersonHandler struct {
    personService services.PersonService
}

func NewPersonHandler(service services.PersonService) *PersonHandler {
    return &PersonHandler{personService: service}
}

// Implementación de los endpoints HTTP...
```

#### Equipment Handlers
```go
// internal/handlers/equipment_handlers.go
package handlers

type EquipmentHandler struct {
    equipmentService services.EquipmentService
}

func NewEquipmentHandler(service services.EquipmentService) *EquipmentHandler {
    return &EquipmentHandler{equipmentService: service}
}

// Implementación de los endpoints HTTP...
```

### 6. Tests Unitarios

#### Person Service Tests
```go
// internal/services/person_service_test.go
package services_test

func TestPersonService_CreatePerson(t *testing.T) {
    // Configuración de test
    mockRepo := new(mocks.PersonRepository)
    service := NewPersonService(mockRepo)
    
    // Caso de prueba
    t.Run("Success", func(t *testing.T) {
        // Configurar expectativas
        expectedPerson := &models.Person{
            FirstName: "John",
            LastName:  "Doe",
            Email:     "john.doe@example.com",
        }
        
        mockRepo.On("Create", mock.AnythingOfType("*models.Person")).
            Return(nil).
            Run(func(args mock.Arguments) {
                p := args.Get(0).(*models.Person)
                p.ID = 1
            })
        
        // Ejecutar
        err := service.RegisterPerson(expectedPerson)
        
        // Verificar
        assert.NoError(t, err)
        assert.Equal(t, uint(1), expectedPerson.ID)
        mockRepo.AssertExpectations(t)
    })
    
    // Más casos de prueba...
}
```

#### Equipment Service Tests
```go
// internal/services/equipment_service_test.go
package services_test

func TestEquipmentService_AssignEquipment(t *testing.T) {
    // Configuración de test
    mockEqRepo := new(mocks.EquipmentRepository)
    mockPersonRepo := new(mocks.PersonRepository)
    service := NewEquipmentService(mockEqRepo, mockPersonRepo)
    
    t.Run("Success", func(t *testing.T) {
        // Configurar expectativas
        equipmentID := uint(1)
        personID := uint(1)
        
        mockPersonRepo.On("FindByID", personID).
            Return(&models.Person{ID: personID}, nil)
            
        mockEqRepo.On("FindByID", equipmentID).
            Return(&models.Equipment{ID: equipmentID, Status: models.StatusAvailable}, nil)
            
        mockEqRepo.On("Update", mock.AnythingOfType("*models.Equipment")).
            Return(nil)
        
        // Ejecutar
        err := service.AssignEquipment(equipmentID, personID)
        
        // Verificar
        assert.NoError(t, err)
        mockPersonRepo.AssertExpectations(t)
        mockEqRepo.AssertExpectations(t)
    })
    
    // Más casos de prueba...
}
```

### 7. Tests de Integración

```go
// integration/person_integration_test.go
package integration

func TestPersonCRUD(t *testing.T) {
    // Configuración de la base de datos de prueba
    db := setupTestDB()
    defer db.Close()
    
    // Inicializar repositorios y servicios
    personRepo := repositories.NewPersonRepository(db)
    personService := services.NewPersonService(personRepo)
    
    // Crear persona
    t.Run("CreatePerson", func(t *testing.T) {
        person := &models.Person{
            FirstName: "Test",
            LastName:  "User",
            Email:     "test@example.com",
        }
        
        err := personService.RegisterPerson(person)
        assert.NoError(t, err)
        assert.NotZero(t, person.ID)
        
        // Verificar que se puede recuperar
        found, err := personService.GetPerson(person.ID)
        assert.NoError(t, err)
        assert.Equal(t, person.Email, found.Email)
    })
    
    // Más pruebas de integración...
}

// integration/equipment_integration_test.go
package integration

func TestEquipmentAssignment(t *testing.T) {
    // Configuración de la base de datos de prueba
    db := setupTestDB()
    defer db.Close()
    
    // Inicializar repositorios y servicios
    personRepo := repositories.NewPersonRepository(db)
    equipmentRepo := repositories.NewEquipmentRepository(db)
    equipmentService := services.NewEquipmentService(equipmentRepo, personRepo)
    
    // Crear persona
    person := &models.Person{
        FirstName: "Test",
        LastName:  "User",
        Email:     "test@example.com",
    }
    personRepo.Create(person)
    
    // Crear equipo
    equipment := &models.Equipment{
        Name:         "Laptop Test",
        Type:         models.EquipmentTypeLaptop,
        Status:       models.StatusAvailable,
        SerialNumber: "TEST123",
    }
    equipmentRepo.Create(equipment)
    
    // Asignar equipo
    t.Run("AssignEquipment", func(t *testing.T) {
        err := equipmentService.AssignEquipment(equipment.ID, person.ID)
        assert.NoError(t, err)
        
        // Verificar asignación
        eq, err := equipmentService.GetEquipment(equipment.ID)
        assert.NoError(t, err)
        assert.Equal(t, models.StatusAssigned, eq.Status)
        assert.Equal(t, person.ID, *eq.AssignedTo)
    })
    
    // Más pruebas de integración...
}
```

### 8. Endpoints de la API

#### Person Endpoints
```
GET    /api/v1/persons          # Listar personas
POST   /api/v1/persons          # Crear persona
GET    /api/v1/persons/:id      # Obtener persona
PUT    /api/v1/persons/:id      # Actualizar persona
DELETE /api/v1/persons/:id      # Eliminar persona
```

#### Equipment Endpoints
```
GET    /api/v1/equipment                     # Listar equipos
POST   /api/v1/equipment                     # Crear equipo
GET    /api/v1/equipment/:id                 # Obtener equipo
PUT    /api/v1/equipment/:id                 # Actualizar equipo
DELETE /api/v1/equipment/:id                 # Eliminar equipo
POST   /api/v1/equipment/:id/assign          # Asignar equipo
POST   /api/v1/equipment/:id/unassign        # Desasignar equipo
GET    /api/v1/equipment/assigned/:personId  # Equipos asignados a una persona
```

### 9. Próximos Pasos

1. **Validaciones**
   - Añadir validaciones de datos de entrada
   - Implementar reglas de negocio específicas

2. **Seguridad**
   - Añadir autenticación/autorización a los endpoints
   - Implementar rate limiting

3. **Documentación**
   - Documentar la API con Swagger/OpenAPI
   - Crear guías de uso

4. **Mejoras**
   - Añadir búsqueda avanzada con filtros
   - Implementar paginación en los listados
   - Añadir historial de cambios

## 📝 Notas de Implementación

### Estructura del Proyecto
- `/backend` - Código fuente del servidor
  - `/internal/handlers` - Manejadores de endpoints HTTP
  - `/internal/services` - Lógica de negocio
  - `/pkg` - Código reusable
  - `/testutils` - Utilidades para pruebas

### Pruebas
- Todas las nuevas características deben incluir:
  - Pruebas unitarias
  - Pruebas de integración
  - Pruebas de extremo a extremo cuando sea aplicable

### Seguridad
- Usar siempre HTTPS en producción
- Validar todas las entradas del usuario
- Implementar rate limiting
- Mantener dependencias actualizadas
