# Pruebas Pendientes - Priorizadas

## 📊 Estado Actual de Cobertura
- `services`: 81.6% ✅
- `services/auth`: 66.7% 🔄 (cubriendo casos básicos, falta renovación de tokens)
- `handlers`: 62.7% 🔄 (cubriendo flujos principales, faltan casos de error específicos)
- `middleware`: 85.2% ✅ (cobertura sólida, probados los casos principales)
- `repositories`: 0% ❌ (faltan pruebas para repositorios)
- `validators`: 0% ❌ (faltan pruebas para validaciones)
- `config`: 0% ❌ (sin pruebas de configuración)
- `database`: 0% ❌ (faltan pruebas de migración y conexión)
- `models`: 0% ❌ (faltan pruebas de modelos)
- `router`: 0% ❌ (sin pruebas de enrutamiento)
- `di`: 0% ❌ (sin pruebas de inyección de dependencias)
- `integration/auth`: 75% 🔄 (cubriendo flujos principales de autenticación)
- `integration/mailer`: 65% 🔄 (pruebas básicas implementadas)


## 🟢 Alta Prioridad (Seguridad y Funcionalidad Crítica)

### 1. Middleware de Autenticación ✅
- [x] Probar autenticación exitosa con token válido
- [x] Manejar token faltante en el header
- [x] Manejar token con formato inválido
- [x] Manejar token expirado
- [x] Validar claims requeridos en el token
- [x] Probar rutas protegidas sin autenticación
- [x] Probar rutas protegidas con token inválido

### 2. Autenticación 🔄
- [x] Validar tokens JWT con firma incorrecta
- [x] Manejar tokens expirados en el middleware (cubierto en `TestValidateToken_Expired`)
- [x] Validar claims requeridos en tokens (cubierto en `TestValidateToken_MissingClaims` y `TestValidateToken_WrongAudience`)
- [x] Probar manejo de JSON inválido en endpoints de autenticación
- [ ] Probar revocación de tokens (requiere implementación)
- [ ] Probar renovación de tokens de refresco
- [ ] Validar expiración de tokens en diferentes escenarios
- [ ] Probar manejo de tokens revocados

### 2. Registro/Login
- [x] Validar fortaleza de contraseñas (cubierto en `TestValidateStrongPassword`)
- [X] Validar formato de emails
- [x] Manejar intentos de registro con email existente (cubierto en `TestRegister_EmailExists`)
- [x] Validar formato de tokens generados (cubierto en múltiples pruebas de validación de tokens)

## 🟡 Media Prioridad (Funcionalidad Importante)

### 1. Servicio de Autenticación
- [ ] Probar renovación de tokens de refresco
- [ ] Validar expiración de tokens
- [ ] Manejar errores del repositorio
- [ ] Probar validación de roles de usuario
- [ ] Probar casos de error en el servicio de autenticación
- [ ] Validar hashing de contraseñas

### 2. Manejadores HTTP
- [x] Validar respuestas HTTP correctas (parcialmente implementado)
- [x] Probar manejo de errores (parcialmente implementado)
- [x] Validar códigos de estado HTTP (parcialmente implementado)
- [ ] Probar validación de entradas
- [ ] Probar límites de tasa (rate limiting)
- [ ] Probar manejo de CORS

## 🔴 Baja Prioridad (Mejoras y Cobertura)

### 1. Pruebas de Integración
- [ ] Probar flujos completos de autenticación
- [ ] Probar interacción con base de datos
- [ ] Probar caché de tokens
- [ ] Probar integración con servicio de correo
- [ ] Probar recuperación de contraseña

### 2. Seguridad Adicional
- [ ] Implementar rate limiting
- [ ] Probar headers de seguridad
- [ ] Validar CORS
- [ ] Probar protección contra CSRF
- [ ] Probar sanitización de entradas

### 3. Rendimiento
- [ ] Probar tiempos de respuesta
- [ ] Probar manejo de carga
- [ ] Probar concurrencia

## 🎯 Objetivos de Cobertura
- [ ] `handlers`: >80% (actual: 62.7%)
- [ ] `services`: >90% (actual: 81.6%)
- [ ] `services/auth`: >85% (actual: 66.7%)
- [ ] `middleware`: >80% (actual: 0%)
- [ ] `repositories`: >80% (actual: 0%)
- [ ] `validators`: >90% (actual: 0%)

## 📋 Próximos Pasos Recomendados
1. Completar pruebas de integración para autenticación (alta prioridad)
2. Mejorar cobertura de `handlers` con más casos de error (alta prioridad)
3. Implementar pruebas para `repositories` (prioridad media)
4. Agregar pruebas para `validators` (prioridad media)
5. Cubrir `config` y `router` (prioridad baja)
6. Implementar pruebas de rendimiento (baja prioridad)

## 📝 Notas Técnicas
- Usar tablas de pruebas para casos similares (ejemplo en `TestAuthHandler_InvalidJSON`)
- Mantener los tests independientes y aislados
- Usar mocks para dependencias externas (ejemplo: `mocks.MockAuthService`)
- Ejecutar `go test -cover ./...` para ver cobertura actual
- Usar `-coverprofile=coverage.out` para generar reporte detallado
- Considerar usar `github.com/stretchr/testify` para aserciones
- Ejecutar pruebas de integración con `go test -tags=integration ./...`
- Mantener la cobertura de pruebas por encima del 80% en componentes críticos
