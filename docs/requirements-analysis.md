# CertiTrack - Análisis de Requisitos del Sistema

## Visión General del Proyecto
CertiTrack es una aplicación web para la gestión integral de certificaciones de personas y equipos dentro de organizaciones pequeñas (50-200 usuarios, ~1000-5000 certificaciones).

## Requisitos de Negocio

### RN001 - Gestión de Certificaciones
- **Alcance**: Operaciones CRUD para certificaciones
- **Detalles**: Registrar, consultar, actualizar, eliminar certificaciones
- **Asociaciones**: Vincular con personas o equipos
- **Datos**: Fecha de emisión, fecha de vencimiento, documentos adjuntos
- **Prioridad**: Alta

### RN002 - Sistema de Alertas y Notificaciones
- **Alcance**: Sistema de notificación proactiva
- **Detalles**: Notificaciones por correo electrónico con avisos configurables (30, 15, 7, 1 días antes del vencimiento)
- **Configuración**: Gestión flexible de tiempos y destinatarios
- **Prioridad**: Alta

### RN003 - Gestión de Personas
- **Alcance**: Gestión de personal
- **Detalles**: Perfiles de empleados, información de contacto, relaciones con certificaciones
- **Prioridad**: Alta

### RN004 - Gestión de Equipos
- **Alcance**: Gestión de equipos/activos
- **Detalles**: Registro de equipos, números de activo, relaciones con certificaciones
- **Prioridad**: Alta

### RN005 - Trazabilidad y Auditoría
- **Alcance**: Sistema de registro de auditoría
- **Detalles**: Registrar acciones clave y cambios en certificaciones
- **Prioridad**: Media

### RN006 - Reportes y Paneles de Control
- **Alcance**: Informes y paneles de control
- **Detalles**: Visualizaciones e informes sobre el estado de las certificaciones
- **Prioridad**: Media

## Restricciones Técnicas

### Restricciones de Escala
- **Usuarios**: 50-200 usuarios concurrentes
- **Volumen de Datos**: 1,000-5,000 certificaciones
- **Crecimiento**: Diseñado para escala de pequeña organización

### Tecnologías Utilizadas
- **Frontend**: React/Next.js
- **Backend**: Go con framework Gin y ORM GORM
- **Base de Datos**: PostgreSQL
- **Despliegue**: Instancia única de AWS EC2 con Docker
- **Enfoque**: Despliegue de MVP optimizado en costos

### Roles de Usuario
- **Administradores**: Acceso completo al sistema, gestión de usuarios, configuración
- **Usuarios Regulares**: Ver sus propias certificaciones, entrada limitada de datos

### Tipos de Certificación
- **Certificaciones de Seguridad**: Primeros auxilios, seguridad contra incendios, seguridad laboral
- **Certificaciones Profesionales**: Licencias, calificaciones técnicas
- **Certificaciones de Equipos**: Calibración, mantenimiento, inspecciones de seguridad
- **Flexibilidad**: Soporte para períodos de renovación variables y tipos personalizados

## Requisitos No Funcionales

### Rendimiento
- **Tiempo de Respuesta**: < 2 segundos para operaciones estándar
- **Usuarios Concurrentes**: Soporte para 50-100 usuarios concurrentes
- **Base de Datos**: Optimizada para cargas de trabajo con muchas lecturas (informes/paneles)

### Seguridad
- **Autenticación**: Sistema seguro de autenticación de usuarios
- **Autorización**: Control de acceso basado en roles
- **Protección de Datos**: Manejo seguro de datos personales y de certificación
- **Registro de Auditoría**: Registro completo de operaciones sensibles

### Fiabilidad
- **Disponibilidad**: Objetivo de 99% de tiempo activo
- **Copia de Seguridad**: Copias de seguridad automáticas diarias
- **Recuperación**: Capacidad de recuperación en un punto específico

### Usabilidad
- **Interfaz**: Interfaz web intuitiva
- **Móvil**: Diseño responsivo para acceso móvil
- **Accesibilidad**: Cumplimiento básico de accesibilidad

### Escalabilidad
- **Horizontal**: Diseñado para despliegue de instancia única inicialmente
- **Vertical**: Capacidad de escalar la instancia EC2 según sea necesario
- **Futuro**: La arquitectura debe admitir despliegue multi-instancia

## Requisitos de Integración

### Sistema de Correo Electrónico
- **Integración SMTP**: Para el envío de notificaciones
- **Plantillas**: Plantillas de correo electrónico configurables
- **Programación**: Programación automática de notificaciones

### Gestión de Archivos
- **Carga de Documentos**: Soporte para documentos de certificación
- **Tipos de Archivo**: PDF, imágenes, formatos de documento comunes
- **Almacenamiento**: Almacenamiento seguro de archivos con controles de acceso

### Informes
- **Formatos de Exportación**: PDF, Excel, CSV
- **Informes Programados**: Generación automática de informes
- **Panel de Control**: Visualización en tiempo real del estado

## Cumplimiento y Regulaciones
- **Privacidad de Datos**: Principios de protección de datos similares al GDPR
- **Requisitos de Auditoría**: Registro de auditoría completo para cumplimiento
- **Retención de Documentos**: Políticas de retención configurables