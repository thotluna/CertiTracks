# CertiTrack - Especificación de la API

## Resumen de la API

La API de CertiTrack sigue los principios RESTful con intercambio de datos en JSON. Todos los endpoints requieren autenticación excepto el endpoint de inicio de sesión.

**URL Base**: `http://localhost:8080/api/v1`

**Autenticación**: Token JWT Bearer en el encabezado Authorization

## Formatos Comunes de Respuesta

### Respuesta Exitosa
```json
{
  "success": true,
  "data": {}, // Datos de respuesta
  "message": "Operación completada exitosamente"
}
```

### Respuesta de Error
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {} // Additional error details
  }
}
```

### Respuesta Paginada
```json
{
  "success": true,
  "data": {
    "items": [], // Array of items
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "totalPages": 5
    }
  }
}
```

## Endpoints de Autenticación

### POST /auth/login
Autentica al usuario y devuelve un token JWT.

**Cuerpo de la Solicitud:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "role": "admin"
    },
    "token": "jwt_token_here",
    "expiresAt": "2024-01-01T12:00:00Z"
  }
}
```

### POST /auth/refresh
Actualizar token JWT.

**Cuerpo de la Solicitud:**
```json
{
  "refreshToken": "refresh_token_here"
}
```

### POST /auth/logout
Cerrar sesión e invalidar el token.

**Encabezados:** `Authorization: Bearer <token>`

## Endpoints de Gestión de Usuarios

### GET /users
Obtener lista de usuarios (solo administradores).

**Parámetros de Consulta:**
- `page` (int): Número de página (por defecto: 1)
- `limit` (int): Elementos por página (por defecto: 20)
- `search` (string): Buscar por nombre o correo electrónico
- `role` (string): Filtrar por rol

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "email": "user@example.com",
        "firstName": "John",
        "lastName": "Doe",
        "role": "admin",
        "isActive": true,
        "createdAt": "2024-01-01T12:00:00Z",
        "lastLogin": "2024-01-01T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "totalPages": 5
    }
  }
}
```

### POST /users
Crear nuevo usuario (solo administradores).

**Cuerpo de la Solicitud:**
```json
{
  "email": "newuser@example.com",
  "password": "password123",
  "firstName": "Jane",
  "lastName": "Smith",
  "phone": "+1234567890",
  "role": "user"
}
```

### GET /users/{id}
Obtener usuario por ID.

### PUT /users/{id}
Actualizar información del usuario.

### DELETE /users/{id}
Eliminar usuario (solo administradores).

## Endpoints de Gestión de Personas

### GET /people
Obtener lista de personas.

**Parámetros de Consulta:**
- `page`, `limit`: Paginación
- `search`: Buscar por nombre, correo electrónico o ID de empleado
- `department`: Filtrar por departamento
- `isActive`: Filtrar por estado activo

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "employeeId": "EMP001",
        "firstName": "John",
        "lastName": "Doe",
        "email": "john.doe@company.com",
        "phone": "+1234567890",
        "department": "Engineering",
        "position": "Software Engineer",
        "hireDate": "2023-01-15",
        "isActive": true,
        "createdAt": "2024-01-01T12:00:00Z",
        "certificationCount": 5,
        "expiringCertifications": 2
      }
    ],
    "pagination": {}
  }
}
```

### POST /people
Crear nueva persona.

**Cuerpo de la Solicitud:**
```json
{
  "employeeId": "EMP002",
  "firstName": "Jane",
  "lastName": "Smith",
  "email": "jane.smith@company.com",
  "phone": "+1234567890",
  "department": "HR",
  "position": "HR Manager",
  "hireDate": "2023-02-01"
}
```

### GET /people/{id}
Obtener persona por ID con información detallada.

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "employeeId": "EMP001",
    "firstName": "John",
    "lastName": "Doe",
    "email": "john.doe@company.com",
    "phone": "+1234567890",
    "department": "Engineering",
    "position": "Software Engineer",
    "hireDate": "2023-01-15",
    "isActive": true,
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:00Z",
    "certifications": [
      {
        "id": "uuid",
        "certificationType": {
          "id": "uuid",
          "name": "First Aid Certification",
          "category": "safety"
        },
        "certificateNumber": "FA-2024-001",
        "issueDate": "2024-01-01",
        "expirationDate": "2026-01-01",
        "status": "active",
        "daysUntilExpiration": 365
      }
    ]
  }
}
```

### PUT /people/{id}
Actualizar información de la persona.

### DELETE /people/{id}
Eliminar persona (eliminación lógica).

## Endpoints de Gestión de Equipos

### GET /equipment
Obtener lista de equipos.

**Parámetros de Consulta:**
- `page`, `limit`: Paginación
- `search`: Buscar por nombre, número de activo o número de serie
- `location`: Filtrar por ubicación
- `isActive`: Filtrar por estado activo

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "assetNumber": "EQ-001",
        "name": "Forklift #1",
        "description": "Electric forklift for warehouse operations",
        "manufacturer": "Toyota",
        "model": "8FBE20",
        "serialNumber": "TY123456",
        "location": "Warehouse A",
        "purchaseDate": "2022-01-15",
        "isActive": true,
        "createdAt": "2024-01-01T12:00:00Z",
        "certificationCount": 3,
        "expiringCertifications": 1
      }
    ],
    "pagination": {}
  }
}
```

### POST /equipment
Crear nuevo equipo.

**Cuerpo de la Solicitud:**
```json
{
  "assetNumber": "EQ-002",
  "name": "Crane #2",
  "description": "Overhead crane for heavy lifting",
  "manufacturer": "Konecranes",
  "model": "CXT",
  "serialNumber": "KC789012",
  "location": "Production Floor",
  "purchaseDate": "2023-03-01"
}
```

### GET /equipment/{id}
Obtener equipo por ID con sus certificaciones.

### PUT /equipment/{id}
Actualizar información del equipo.

### DELETE /equipment/{id}
Eliminar equipo (eliminación lógica).

## Endpoints de Tipos de Certificación

### GET /certification-types
Obtener lista de tipos de certificación.

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "name": "First Aid Certification",
        "description": "Basic first aid and CPR certification",
        "category": "safety",
        "defaultValidityPeriod": 730,
        "requiresRenewal": true,
        "isActive": true,
        "createdAt": "2024-01-01T12:00:00Z"
      }
    ]
  }
}
```

### POST /certification-types
Crear nuevo tipo de certificación (solo administradores).

**Cuerpo de la Solicitud:**
```json
{
  "name": "Forklift Operation License",
  "description": "License to operate forklift equipment",
  "category": "professional",
  "defaultValidityPeriod": 1095,
  "requiresRenewal": true
}
```

### GET /certification-types/{id}
Obtener tipo de certificación por ID.

### PUT /certification-types/{id}
Actualizar tipo de certificación (solo administradores).

### DELETE /certification-types/{id}
Eliminar tipo de certificación (solo administradores).

## Endpoints de Certificaciones

### GET /certifications
Obtener lista de certificaciones.

**Parámetros de Consulta:**
- `page`, `limit`: Paginación
- `search`: Buscar por número de certificado o autoridad emisora
- `personId`: Filtrar por persona
- `equipmentId`: Filtrar por equipo
- `certificationTypeId`: Filtrar por tipo de certificación
- `status`: Filtrar por estado (active, expired, revoked, pending)
- `expiringInDays`: Filtrar certificaciones que expiran en X días
- `sortBy`: Ordenar por campo (expirationDate, issueDate, createdAt)
- `sortOrder`: Orden de clasificación (asc, desc)

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "certificationType": {
          "id": "uuid",
          "name": "First Aid Certification",
          "category": "safety"
        },
        "person": {
          "id": "uuid",
          "firstName": "John",
          "lastName": "Doe",
          "employeeId": "EMP001"
        },
        "equipment": null,
        "certificateNumber": "FA-2024-001",
        "issuingAuthority": "Red Cross",
        "issueDate": "2024-01-01",
        "expirationDate": "2026-01-01",
        "status": "active",
        "notes": "Completed advanced first aid course",
        "daysUntilExpiration": 365,
        "documentCount": 2,
        "createdAt": "2024-01-01T12:00:00Z"
      }
    ],
    "pagination": {}
  }
}
```

### POST /certifications
Crear nueva certificación.

**Cuerpo de la Solicitud:**
```json
{
  "certificationTypeId": "uuid",
  "personId": "uuid", // Either personId or equipmentId, not both
  "equipmentId": null,
  "certificateNumber": "FA-2024-002",
  "issuingAuthority": "Red Cross",
  "issueDate": "2024-01-15",
  "expirationDate": "2026-01-15",
  "notes": "Initial certification"
}
```

### GET /certifications/{id}
Obtener certificación por ID con detalles completos.

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "certificationType": {
      "id": "uuid",
      "name": "First Aid Certification",
      "description": "Basic first aid and CPR certification",
      "category": "safety"
    },
    "person": {
      "id": "uuid",
      "firstName": "John",
      "lastName": "Doe",
      "employeeId": "EMP001",
      "email": "john.doe@company.com"
    },
    "equipment": null,
    "certificateNumber": "FA-2024-001",
    "issuingAuthority": "Red Cross",
    "issueDate": "2024-01-01",
    "expirationDate": "2026-01-01",
    "status": "active",
    "notes": "Completed advanced first aid course",
    "documents": [
      {
        "id": "uuid",
        "fileName": "certificate.pdf",
        "fileSize": 1024000,
        "uploadedAt": "2024-01-01T12:00:00Z"
      }
    ],
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:00Z"
  }
}
```

### PUT /certifications/{id}
Actualizar certificación.

### DELETE /certifications/{id}
Eliminar certificación.

## Document Management Endpoints

### POST /certifications/{id}/documents
Subir documento para certificación.

**Content-Type**: `multipart/form-data`

**Form Data:**
- `file`: Archivo a subir
- `description`: Descripción opcional

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "fileName": "certificate.pdf",
    "fileSize": 1024000,
    "mimeType": "application/pdf",
    "uploadedAt": "2024-01-01T12:00:00Z"
  }
}
```

### GET /certifications/{certId}/documents/{docId}
Descargar documento.

**Respuesta**: Descarga del archivo con los encabezados apropiados.

### DELETE /certifications/{certId}/documents/{docId}
Eliminar documento.

## Notification Endpoints

### GET /notifications
Obtener notificaciones para el usuario actual.

**Parámetros de Consulta:**
- `page`, `limit`: Paginación
- `status`: Filtrar por estado (pending, sent, failed)
- `certificationId`: Filtrar por certificación

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "certification": {
          "id": "uuid",
          "certificateNumber": "FA-2024-001",
          "certificationType": {
            "name": "First Aid Certification"
          },
          "person": {
            "firstName": "John",
            "lastName": "Doe"
          }
        },
        "subject": "Certification Expiring Soon",
        "message": "Your First Aid Certification will expire in 30 days",
        "status": "sent",
        "scheduledFor": "2024-01-01T09:00:00Z",
        "sentAt": "2024-01-01T09:00:00Z"
      }
    ],
    "pagination": {}
  }
}
```

### POST /notifications/test
Enviar notificación de prueba (solo administradores).

**Cuerpo de la Solicitud:**
```json
{
  "recipientEmail": "test@example.com",
  "subject": "Test Notification",
  "message": "This is a test notification"
}
```

## Reporting Endpoints

### GET /reports/dashboard
Obtener estadísticas del dashboard.

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "totalCertifications": 1250,
    "activeCertifications": 1100,
    "expiredCertifications": 150,
    "expiringIn30Days": 45,
    "expiringIn7Days": 12,
    "totalPeople": 200,
    "totalEquipment": 85,
    "certificationsByCategory": {
      "safety": 600,
      "professional": 400,
      "equipment": 250
    },
    "expirationTrend": [
      {
        "month": "2024-01",
        "expiring": 25,
        "renewed": 20
      }
    ]
  }
}
```

### GET /reports/certifications
Obtener informe de certificaciones con filtros.

**Parámetros de Consulta:**
- `startDate`, `endDate`: Rango de fechas
- `category`: Categoría de certificación
- `status`: Estado de certificación
- `format`: Export format (json, csv, pdf)

### GET /reports/people/{id}/certifications
Obtener informe de certificaciones para una persona específica.

### GET /reports/equipment/{id}/certifications
Obtener informe de certificaciones para un equipo específico.

## Health Check Endpoints

### GET /health
Verificar estado del sistema.

**Respuesta:**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z",
    "version": "1.0.0",
    "services": {
      "database": "healthy",
      "redis": "healthy",
      "email": "healthy"
    }
  }
}
```

## Error Codes

| Código | Descripción |
|------|-------------|
| `INVALID_REQUEST` | Petición inválida |
| `UNAUTHORIZED` | Autenticación requerida |
| `FORBIDDEN` | Insufficient permissions |
| `NOT_FOUND` | Recurso no encontrado |
| `CONFLICT` | Recurso ya existe |
| `VALIDATION_ERROR` | Validación de datos fallida |
| `INTERNAL_ERROR` | Error interno del servidor |
| `SERVICE_UNAVAILABLE` | Servicio no disponible |

## Rate Limiting

- **General API**: 1000 peticiones por hora por usuario
- **File Upload**: 100 peticiones por hora por usuario
- **Authentication**: 10 peticiones por minuto por IP

## API Versioning

La API utiliza la versión URL (`/api/v1/`). Las versiones futuras se introducirán según sea necesario manteniendo la compatibilidad hacia atrás para al menos una versión mayor.