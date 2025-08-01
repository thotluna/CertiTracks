# Pruebas Pendientes - Priorizadas

## 📊 Estado Actual de Cobertura
- `services`: 81.6% ✅
- `services/auth`: 66.7% 🔄
- `handlers`: 62.7% 🔄
- `middleware`: 85.2% ✅
- `repositories`: 0% ❌
- `validators`: 0% ❌
- `config`: 0% ❌
- `database`: 0% ❌
- `models`: 0% ❌
- `router`: 0% ❌
- `di`: 0% ❌


## 🟢 Alta Prioridad (Seguridad y Funcionalidad Crítica)

### 1. Middleware de Autenticación ✅
- [x] Probar autenticación exitosa con token válido
- [x] Manejar token faltante en el header
- [x] Manejar token con formato inválido
- [x] Manejar token expirado
- [x] Validar claims requeridos en el token
- [x] Probar rutas protegidas sin autenticación
- [x] Probar rutas protegidas con token inválido

### 2. Autenticación ✅
- [x] Validar tokens JWT con firma incorrecta
- [x] Manejar tokens expirados en el middleware (cubierto en `TestValidateToken_Expired`)
- [x] Validar claims requeridos en tokens (cubierto en `TestValidateToken_MissingClaims` y `TestValidateToken_WrongAudience`)
- [ ] Probar revocación de tokens (requiere implementación)

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

### 2. Manejadores HTTP
- [ ] Validar respuestas HTTP correctas
- [ ] Probar manejo de errores
- [ ] Validar códigos de estado HTTP
- [ ] Probar validación de entradas

## 🔴 Baja Prioridad (Mejoras y Cobertura)

### 1. Pruebas de Integración
- [ ] Probar flujos completos de autenticación
- [ ] Probar interacción con base de datos
- [ ] Probar caché de tokens

### 2. Seguridad Adicional
- [ ] Implementar rate limiting
- [ ] Probar headers de seguridad
- [ ] Validar CORS

## 🎯 Objetivos de Cobertura
- [ ] `handlers`: >80% (actual: 62.7%)
- [ ] `services`: >90% (actual: 81.6%)
- [ ] `services/auth`: >85% (actual: 66.7%)
- [ ] `middleware`: >80% (actual: 0%)
- [ ] `repositories`: >80% (actual: 0%)
- [ ] `validators`: >90% (actual: 0%)

## 📋 Próximos Pasos Recomendados
1. Agregar pruebas para `middleware` (prioridad alta)
2. Mejorar cobertura de `handlers` (prioridad alta)
3. Agregar pruebas para `validators` (prioridad media)
4. Implementar pruebas para `repositories` (prioridad media)
5. Cubrir `config` y `router` (prioridad baja)

## 📝 Notas Técnicas
- Usar tablas de pruebas para casos similares
- Mantener los tests independientes
- Usar mocks para dependencias externas
- Ejecutar `go test -cover ./...` para ver cobertura actual
