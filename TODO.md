# Pendientes de Autenticación

## Implementación de Recuperación de Contraseña (TDD)

### Fase 1: Configuración Inicial
1. [x] Crear modelo `PasswordResetToken`
2. [ ] Configurar servicio de correo (ej: SMTP o servicio externo)

### Fase 2: Solicitud de Restablecimiento
1. [ ] Test: Endpoint POST `/auth/forgot-password`
   - Validar correo requerido
   - Validar formato de correo
   - No revelar si el correo existe
2. [ ] Test: Generación de token seguro
3. [ ] Test: Guardar token en base de datos
4. [ ] Test: Envío de correo con enlace
5. [ ] Implementar lógica del endpoint

### Fase 3: Validación de Token
1. [ ] Test: Endpoint GET `/auth/validate-reset-token/{token}`
   - Token válido: 200 OK
   - Token inválido: 400 Bad Request
   - Token expirado: 400 Bad Request
   - Token ya usado: 400 Bad Request
2. [ ] Implementar validación de token

### Fase 4: Restablecer Contraseña
1. [ ] Test: Endpoint POST `/auth/reset-password`
   - Validar token requerido
   - Validar nueva contraseña (mínimo 8 caracteres, etc.)
   - Confirmación de contraseña
2. [ ] Test: Actualización de contraseña exitosa
   - Hash de la nueva contraseña
   - Marcar token como usado
   - Invalidar tokens anteriores del usuario
3. [ ] Test: Notificación por correo de cambio de contraseña
4. [ ] Implementar lógica de restablecimiento

### Fase 5: Seguridad Adicional
1. [ ] Test: Límite de intentos por IP
2. [ ] Test: Tasa límite para envío de correos
3. [ ] Test: Expiración de tokens (1 hora)
4. [ ] Test: No permitir reutilizar contraseñas recientes

### Fase 6: Documentación
1. [ ] Documentar endpoints en Swagger/OpenAPI
2. [ ] Crear guía de implementación para frontend
3. [ ] Documentar políticas de seguridad

## Prioridad Alta
- [ ] Mejorar gestión de tokens
  - [X] Revocar refresh tokens al hacer logout
  - [ ] Implementar limpieza periódica de tokens expirados
  - [ ] Añadir revocación masiva por usuario (útil en cambio de contraseña)
  - [x] Agregar logs de eventos de revocación
- [ ] Implementar recuperación de contraseña
  - Endpoint para solicitar restablecimiento
  - Envío de correo con enlace seguro
  - Formulario para restablecer contraseña

- [ ] Añadir verificación de correo electrónico
  - Envío de correo de verificación al registrarse
  - Endpoint para validar token de verificación
  - Actualizar estado de verificación del usuario

## Prioridad Media
- [ ] Mejorar seguridad
  - Límite de intentos de inicio de sesión
  - Bloqueo temporal de cuentas tras múltiples intentos fallidos
  - Registro de actividades sospechosas

- [ ] Mejorar usabilidad
  - Información de sesiones activas
  - Cierre de sesión remota
  - Soporte para múltiples dispositivos

## Prioridad Baja
- [ ] Documentación adicional
  - Ejemplos de código para diferentes lenguajes
  - Guía de migración
  - Preguntas frecuentes

## En Progreso
- [x] Sistema básico de autenticación con JWT
- [x] Middleware de autenticación
- [x] Pruebas de integración básicas
