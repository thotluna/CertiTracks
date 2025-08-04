# Pendientes de Autenticación

## Prioridad Alta
- [ ] Mejorar gestión de tokens
  - [ ] Revocar refresh tokens al hacer logout
  - [ ] Implementar limpieza periódica de tokens expirados
  - [ ] Añadir revocación masiva por usuario (útil en cambio de contraseña)
  - [ ] Agregar logs de eventos de revocación
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
