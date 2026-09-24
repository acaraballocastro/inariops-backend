# Capacidades del backend InariOps

> Estado revisado sobre la rama `itinerary` del backend `inariops-backend`.

Este documento describe exclusivamente la funcionalidad actualmente implementada en el backend Go. No se contabilizan funcionalidades existentes únicamente en el frontend.

---

# 1. Estado general

El backend dispone actualmente de una API REST bajo `/api/v1`, organizada alrededor de:

- Go.
- PostgreSQL.
- Gorilla Mux.
- Arquitectura por módulos.
- Separación entre API, application services y módulos de dominio.
- Repositories para persistencia.
- DTOs para entrada/salida.
- Servicios de aplicación para operaciones complejas.
- Middleware y sistema común de respuestas/errores.
- Logging.
- Workers para procesos automáticos.

La rama `itinerary` contiene actualmente las siguientes áreas funcionales:

- Autenticación.
- Usuarios.
- Guías.
- Disponibilidad de guías.
- Clientes.
- Agencias.
- Idiomas.
- Zonas.
- Reservas.
- Días de tour.
- Asignación de guías.
- Firma.
- Voucher.
- Actividades.
- Places.
- Itinerarios.
- Sincronización de estados.
- Workers automáticos.

---

# 2. Resumen de capacidades

| # | Capacidad | Estado |
|---|---|---|
| 1 | Infraestructura HTTP / API | ✅ Implementado |
| 2 | Autenticación | 🟡 Parcial |
| 3 | Usuarios | ✅ Implementado |
| 4 | Guías | ✅ Implementado |
| 5 | Disponibilidad de guías | ✅ Implementado |
| 6 | Clientes | ✅ Implementado |
| 7 | Agencias | ✅ Implementado |
| 8 | Idiomas | ✅ Implementado |
| 9 | Zonas | ✅ Implementado |
| 10 | Reservas | ✅ Implementado |
| 11 | Días de tour | ✅ Implementado |
| 12 | Asignación de guías | ✅ Implementado |
| 13 | Estados de reserva/tour | 🟡 Parcial |
| 14 | Firma | 🟡 Parcial |
| 15 | Voucher | 🟡 Parcial |
| 16 | Actividades | ✅ Implementado |
| 17 | Places | ✅ Implementado |
| 18 | Itinerarios | ✅ Implementado |
| 19 | Workers automáticos | 🟡 Parcial |
| 20 | Historial / auditoría | 🟡 Parcial |
| 21 | Incidencias | ❌ Pendiente |
| 22 | Notificaciones | ❌ Pendiente |
| 23 | Archivos / documentos | ❌ Pendiente |
| 24 | Email | ❌ Pendiente |
| 25 | RBAC / permisos | ❌ Pendiente |
| 26 | Búsqueda avanzada / filtros | 🟡 Parcial |
| 27 | Paginación | ❌ Pendiente |
| 28 | Testing | 🟡 Parcial |
| 29 | Observabilidad | 🟡 Parcial |
| 30 | Integraciones externas | ❌ Pendiente |

---

# 3. Infraestructura HTTP / API

## Estado: ✅ Implementado

El router principal está definido en:

`internal/api/routes.go`

La API utiliza:

```text
/api/v1
```

Actualmente están registradas las siguientes áreas:

```text
/auth
/users
/reservations
/tour-days
/activities
/itinerary
/places
/guides
/customers
/agencies
/languages
/zones
```

También existe:

```http
GET /api/v1/health
```

Se dispone de:

- CORS.
- OPTIONS / preflight.
- Respuestas JSON.
- Handler global para 404.
- Handler global para 405.
- Logging.

---

# 4. Autenticación

## Estado: 🟡 Parcial

Actualmente existe:

```http
POST  /api/v1/auth/
PATCH /api/v1/auth/
```

El login:

- Busca el usuario por email.
- Comprueba que esté activo.
- Valida contraseña mediante bcrypt.
- Genera JWT.
- Devuelve información básica del usuario.
- Informa de si debe cambiar la contraseña.

También existe cambio de contraseña con validación de contraseña anterior y reglas de complejidad.

### Pendiente

- Middleware global de autenticación JWT.
- Protección real de las rutas.
- Autorización por rol.
- Logout.
- Recuperación de contraseña.
- Reset de contraseña.
- Gestión segura/configurable del secreto JWT.

Por tanto, **existe autenticación pero todavía no existe un sistema completo de autorización de la API**.

---

# 5. Usuarios

## Estado: ✅ Implementado

Endpoints:

```http
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/guides
GET    /api/v1/users/{id}
PATCH  /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

Soporta:

- Creación de usuarios.
- Roles `ADMIN` y `GUIDE`.
- Credenciales.
- Actualización.
- Desactivación.
- Creación automática del registro de guía cuando corresponde.

El DELETE implementado para usuarios es una desactivación lógica.

---

# 6. Guías

## Estado: ✅ Implementado

Endpoints:

```http
GET    /api/v1/guides
POST   /api/v1/guides
GET    /api/v1/guides/{id}
PATCH  /api/v1/guides/{id}
GET    /api/v1/guides/user/{user_id}

GET    /api/v1/guides/{id}/languages
POST   /api/v1/guides/{id}/languages
DELETE /api/v1/guides/{id}/languages/{language_id}
```

La creación de una guía contempla:

- Usuario.
- Perfil de guía.
- Idiomas.
- Zonas.
- Información relacionada.

La actualización permite sincronizar asociaciones de idiomas y zonas.

---

# 7. Disponibilidad de guías

## Estado: ✅ Implementado

La rama `itinerary` incorpora funcionalidad específica para disponibilidad.

Incluye:

- Creación de disponibilidades.
- Consulta.
- Actualización.
- Eliminación.
- Detección de conflictos.
- Comprobación de disponibilidad durante la asignación.

También existe lógica para determinar los días de tour disponibles para un guía.

---

# 8. Clientes

## Estado: ✅ Implementado

Endpoints:

```http
GET    /api/v1/customers
POST   /api/v1/customers
GET    /api/v1/customers/{id}
PATCH  /api/v1/customers/{id}
DELETE /api/v1/customers/{id}
```

El backend también dispone de lógica de búsqueda por:

- Nombre.
- Documento.
- Teléfono.
- Email.

Los clientes pueden crearse/reutilizarse automáticamente durante la creación de una reserva.

---

# 9. Agencias

## Estado: ✅ Implementado

Endpoints:

```http
GET    /api/v1/agencies
POST   /api/v1/agencies
GET    /api/v1/agencies/{id}
PATCH  /api/v1/agencies/{id}
DELETE /api/v1/agencies/{id}
```

La lista permite filtros básicos como:

- Nombre.
- Email.
- Estado activo.

Las reservas validan la existencia de la agencia asociada.

---

# 10. Idiomas

## Estado: ✅ Implementado

Endpoints:

```http
GET    /api/v1/languages
POST   /api/v1/languages
GET    /api/v1/languages/{code}
PATCH  /api/v1/languages/{code}
DELETE /api/v1/languages/{code}
```

Existe:

- CRUD.
- Validación de códigos duplicados.
- Asociación con guías.

---

# 11. Zonas

## Estado: ✅ Implementado

Endpoints:

```http
GET    /api/v1/zones
POST   /api/v1/zones
GET    /api/v1/zones/{name}
PATCH  /api/v1/zones/{id}
DELETE /api/v1/zones/{name}
```

Las zonas también se utilizan dentro de la lógica de guías y disponibilidad.

---

# 12. Reservas

## Estado: ✅ Implementado

Endpoints:

```http
GET    /api/v1/reservations
POST   /api/v1/reservations

GET    /api/v1/reservations/{code}
GET    /api/v1/reservations/details/{code}

PATCH  /api/v1/reservations/{code}

PATCH  /api/v1/reservations/{code}/signature
PATCH  /api/v1/reservations/{code}/voucher

DELETE /api/v1/reservations/{code}
DELETE /api/v1/reservations/erase/{code}

GET    /api/v1/reservations/{code}/customers
```

La creación de una reserva realiza varias operaciones relacionadas:

1. Validación de agencia.
2. Generación de ID.
3. Generación de código.
4. Inicialización del estado.
5. Inicialización del estado de firma.
6. Inicialización del voucher.
7. Creación automática de los días de tour.
8. Creación/reutilización de clientes.
9. Creación de relaciones reserva-cliente.

La actualización permite modificar parcialmente la reserva y sincronizar los clientes asociados.

También existe:

- Cancelación lógica.
- Eliminación física.
- Actualización independiente de firma.
- Actualización independiente de voucher.

---

# 13. Días de tour

## Estado: ✅ Implementado

Endpoints:

```http
POST   /api/v1/tour-days
GET    /api/v1/tour-days/{id}
PATCH  /api/v1/tour-days/{id}
DELETE /api/v1/tour-days/{id}

GET    /api/v1/tour-days/by-reservation/{reservation_id}

GET    /api/v1/tour-days/available-for-guide/{guide_id}

PATCH  /api/v1/tour-days/assign-guide
PATCH  /api/v1/tour-days/unassign-guide

PATCH  /api/v1/tour-days/{id}/confirm
```

Un día de tour soporta información como:

- Fecha.
- Hora.
- Duración.
- Zona.
- Guía.
- Personas.
- Punto de encuentro.
- Hoteles.
- Notas.
- Remuneración.
- Estado.
- Voucher.

Además, el backend permite trabajar con múltiples días por reserva.

---

# 14. Asignación de guías

## Estado: ✅ Implementado

Actualmente existe una capa de aplicación específica para operaciones de tour days.

La asignación permite:

- Asignar uno o varios días.
- Comprobar disponibilidad.
- Comprobar límites diarios.
- Preasignar guía.
- Confirmar posteriormente.
- Desasignar.
- Registrar historial.

La disponibilidad del guía está integrada en el proceso.

---

# 15. Estados de reservas y días

## Estado: 🟡 Parcial

Existe una lógica de estados relativamente completa.

Se manejan estados como:

```text
PENDING
PENDING_ASSIGNMENT
GUIDE_PREASSIGNED
PAYMENT_PENDING
PARTIALLY_CONFIRMED
CONFIRMED
COMPLETED
CANCELLED
FORCE_MAJEURE_CANCELLED
```

También existen estados independientes para:

### Firma

```text
NOT_SENT
SENT
SIGNED
REJECTED
```

### Voucher

```text
NOT_GENERATED
GENERATED
PARTIALLY_SENT
SENT
```

También existe un resolver para determinar el estado global de una reserva a partir de sus días.

### Pendiente

La máquina de estados todavía necesita terminar de conectarse con todos los procesos de negocio.

---

# 16. Firma

## Estado: 🟡 Parcial

Actualmente existe:

```http
PATCH /api/v1/reservations/{code}/signature
```

Permite modificar el estado de firma de la reserva.

### Pendiente

- Generación/envío real de documentos.
- Integración con proveedor externo de firma.
- Callback/webhook.
- Seguimiento completo del proceso.
- Gestión de documentos firmados.

---

# 17. Voucher

## Estado: 🟡 Parcial

Existe:

```http
PATCH /api/v1/reservations/{code}/voucher
```

y un estado específico de voucher.

Por tanto, el **workflow de estado está iniciado**.

### Pendiente

- Generación del PDF.
- Plantilla definitiva.
- Almacenamiento del documento.
- Envío por email.
- Tracking completo de entrega.
- Integración con el sistema externo que corresponda.

---

# 18. Actividades

## Estado: ✅ Implementado

Esta es una de las funcionalidades incorporadas recientemente.

Endpoints:

```http
GET    /api/v1/activities
GET    /api/v1/activities/{id}
POST   /api/v1/activities
PATCH  /api/v1/activities/{id}
DELETE /api/v1/activities/{id}
```

Existe una implementación separada con:

- Handler.
- Service.
- Repository.
- DTOs.
- Persistencia.

Por tanto, **Activities ya forma parte del backend funcional de la rama `itinerary`**.

---

# 19. Places

## Estado: ✅ Implementado

Los places están relacionados directamente con un día de tour.

Endpoints:

```http
GET    /api/v1/tour-days/{tour_day_id}/places
POST   /api/v1/tour-days/{tour_day_id}/places

PATCH  /api/v1/tour-days/{tour_day_id}/places/{place_id}
DELETE /api/v1/tour-days/{tour_day_id}/places/{place_id}
```

La implementación se encuentra separada dentro de la capa de aplicación de `tour_days`.

Esto permite asociar lugares concretos a cada día del tour.

---

# 20. Itinerarios

## Estado: ✅ Implementado

Los itinerarios también están vinculados a un día de tour.

Endpoints:

```http
GET    /api/v1/tour-days/{tour_day_id}/itinerary
POST   /api/v1/tour-days/{tour_day_id}/itinerary

PATCH  /api/v1/tour-days/{tour_day_id}/itinerary/{item_id}
DELETE /api/v1/tour-days/{tour_day_id}/itinerary/{item_id}
```

La implementación incluye una capa específica de aplicación para itinerary.

Por tanto, el backend ya soporta:

```text
Tour Day
 ├── Places
 └── Itinerary
      └── Itinerary Items
```

y Activities existe como catálogo independiente que puede utilizarse dentro de la lógica de itinerarios.

---

# 21. Workers automáticos

## Estado: 🟡 Parcial

Existen workers ejecutados desde el backend para procesos automáticos.

Entre ellos:

### Expiración de asignaciones

Comprueba periódicamente asignaciones preasignadas que hayan superado el límite temporal establecido.

### Sincronización de estados

Recalcula el estado de las reservas a partir de sus días.

Esto significa que ya existe una primera capa de automatización backend.

### Pendiente

- Mayor cobertura de procesos.
- Sistema de jobs más robusto.
- Retries.
- Observabilidad específica.
- Gestión de errores persistentes.

---

# 22. Historial / auditoría

## Estado: 🟡 Parcial

Existe historial asociado a operaciones importantes, especialmente en:

- Asignación de guías.
- Desasignación.
- Expiración de asignaciones.
- Estados de tour days.

El sistema registra información como:

- Actor.
- Acción.
- Razón.
- Momento de ejecución.

### Pendiente

Convertirlo en un sistema de auditoría transversal para todas las entidades y operaciones administrativas.

---

# 23. Incidencias

## Estado: ❌ Pendiente

No existe actualmente un módulo HTTP completo para:

```text
Incidents
```

El modelo funcional previsto contempla incidencias de los guías, pero la implementación backend todavía debe realizarse.

Pendiente:

- Modelo.
- Repository.
- Service.
- Handler.
- Estados.
- Asociación con tour day/reserva.
- API.
- Historial.
- Notificaciones.

---

# 24. Notificaciones

## Estado: ❌ Pendiente

No existe todavía un sistema completo de notificaciones.

Pendiente:

- Notificaciones por zona.
- Notificaciones a guías.
- Eventos de asignación.
- Avisos de expiración.
- Notificaciones de incidencias.
- Email.
- Posibles canales adicionales.

---

# 25. Archivos y documentos

## Estado: ❌ Pendiente

Actualmente no existe un sistema completo de gestión documental.

Pendiente:

- Upload.
- Storage.
- Metadata.
- Asociación con reservas.
- Asociación con vouchers.
- Documentos firmados.
- Descarga.
- Eliminación.

---

# 26. Email

## Estado: ❌ Pendiente

No existe todavía una capa completa para envío de emails.

Especialmente pendiente:

- Envío de vouchers.
- Confirmaciones.
- Avisos a guías.
- Notificaciones.
- Recuperación de contraseña.

---

# 27. RBAC / permisos

## Estado: ❌ Pendiente

Aunque existen roles:

```text
ADMIN
GUIDE
```

actualmente no existe una capa global que fuerce permisos sobre los endpoints.

El JWT se genera, pero las rutas no están protegidas globalmente mediante middleware JWT.

Pendiente:

```text
Authentication
        ↓
JWT Middleware
        ↓
User
        ↓
Role
        ↓
Permission
        ↓
Endpoint
```

---

# 28. Búsqueda y filtros

## Estado: 🟡 Parcial

Ya existen filtros/búsquedas en determinadas áreas.

Ejemplos:

- Clientes por distintos campos.
- Agencias por nombre/email/estado.
- Disponibilidad de guías.
- Días disponibles para un guía.

Sin embargo, todavía no existe un sistema homogéneo de filtros para toda la API.

Pendiente especialmente:

- Reservas por estado.
- Reservas por fechas.
- Reservas por agencia.
- Reservas por guía.
- Tour days por fecha.
- Tour days por zona.
- Búsqueda global.

---

# 29. Paginación

## Estado: ❌ Pendiente

No existe actualmente una estrategia general de paginación para los listados.

Debería incorporarse especialmente en:

```text
/users
/guides
/customers
/reservations
/activities
```

y posteriormente en cualquier colección que pueda crecer significativamente.

---

# 30. Testing

## Estado: 🟡 Parcial

La arquitectura actual facilita la introducción de tests debido a la separación:

```text
Handler
   ↓
Service
   ↓
Repository
```

pero el backend todavía necesita aumentar significativamente su cobertura de pruebas.

Prioridad:

1. Reservations.
2. Tour days.
3. Guide assignment.
4. Availability.
5. Itinerary.
6. Places.
7. Activities.
8. Authentication.

---

# 31. Observabilidad

## Estado: 🟡 Parcial

Existe:

- Logging.
- Health endpoint.
- Logging de errores.
- Sistema común de respuestas.
- Errores de dominio/aplicación.

Pero todavía no existe una infraestructura completa de observabilidad.

Pendiente:

- Métricas.
- Tracing.
- Métricas de negocio.
- Dashboards.
- Alertas.

---

# 32. Integraciones externas

## Estado: ❌ Pendiente

El backend todavía no tiene implementadas completamente las integraciones externas previstas para:

- Firma digital.
- Email.
- Almacenamiento documental.
- Otros proveedores externos.

La arquitectura actual permite incorporarlas posteriormente sin modificar el dominio principal.

---

# 33. Arquitectura actual

La estructura de la rama `itinerary` muestra una evolución hacia una arquitectura modular.

Actualmente encontramos:

```text
internal/
├── api/
│   ├── activity_routes.go
│   ├── agency_routes.go
│   ├── auth_routes.go
│   ├── customer_routes.go
│   ├── guides_routes.go
│   ├── itinerary_routes.go
│   ├── language_routes.go
│   ├── places_route.go
│   ├── reservations_routes.go
│   ├── routes.go
│   ├── tourday_routes.go
│   ├── user_routes.go
│   └── zone_routes.go
│
├── application/
│   ├── guides/
│   ├── reservation_customers/
│   ├── reservations/
│   └── tour_days/
│       ├── itinerary/
│       └── places/
│
├── modules/
│   ├── auth/
│   ├── customers/
│   ├── guides/
│   ├── itinerary/
│   ├── tours/
│   └── users/
│
├── db/
└── shared/
    ├── errors/
    ├── logger/
    ├── middleware/
    └── response/
```

Esto indica que el backend ya no es simplemente un CRUD monolítico, sino que tiene una separación clara entre:

```text
HTTP
 ↓
Application
 ↓
Domain / Modules
 ↓
Repository
 ↓
PostgreSQL
```

---

# 34. Situación actual del proyecto

El backend ya cubre una parte importante del núcleo operativo de InariOps.

Especialmente desarrolladas están estas áreas:

```text
Users
Guides
Guide availability
Customers
Agencies
Reservations
Tour Days
Guide assignment
Activities
Places
Itinerary
```

La incorporación de:

```text
Activities
Places
Itinerary
```

es un cambio importante respecto a versiones anteriores del backend.

El núcleo de planificación de un tour actualmente puede representarse como:

```text
Reservation
     │
     ├── Customers
     │
     └── Tour Days
            │
            ├── Guide
            ├── Availability
            ├── Places
            │
            └── Itinerary
                   │
                   └── Activities
```

---

# 35. Próximas áreas de desarrollo

A partir del estado real actual, las siguientes áreas pendientes tienen especial relevancia:

## Fase 1 — Cerrar el núcleo operativo

- Terminar máquina de estados.
- Revisar todas las transiciones de Reservation/Tour Day.
- Completar reglas de asignación.
- Completar disponibilidad.
- Añadir incidencias.

## Fase 2 — Seguridad

- JWT middleware.
- Protección de endpoints.
- RBAC.
- Permisos por rol.
- Gestión segura del secreto.
- Recuperación de contraseña.

## Fase 3 — Documentos

- Generación de voucher.
- PDF.
- Storage.
- Firma digital.
- Email.

## Fase 4 — Comunicación

- Sistema de notificaciones.
- Email.
- Avisos de asignación.
- Avisos de incidencias.
- Expiraciones.

## Fase 5 — Calidad de plataforma

- Paginación.
- Filtros.
- Tests.
- Métricas.
- Observabilidad.
- Auditoría completa.

---

# 36. Conclusión

El backend actual ya dispone de un **núcleo funcional de gestión de tours**, incluyendo reservas, clientes, guías, disponibilidad, asignaciones y, en la rama `itinerary`, la gestión de actividades, lugares e itinerarios.

Las principales carencias actuales no están ya en el CRUD básico del negocio, sino en las capas que convierten el backend en una plataforma operativa completa:

```text
                 INARIOPS BACKEND

              ┌───────────────────┐
              │   Authentication  │
              │      / RBAC       │
              └─────────┬─────────┘
                        │
              ┌─────────▼─────────┐
              │   CORE OPERATIONS │
              │                   │
              │ Reservations      │
              │ Tour Days         │
              │ Guides            │
              │ Customers         │
              │ Agencies          │
              └─────────┬─────────┘
                        │
              ┌─────────▼─────────┐
              │   TOUR PLANNING   │
              │                   │
              │ Activities        │
              │ Places            │
              │ Itinerary         │
              │ Availability      │
              └─────────┬─────────┘
                        │
          ┌─────────────┼─────────────┐
          │             │             │
     Incidents      Documents    Notifications
          │             │             │
          └─────────────┼─────────────┘
                        │
              ┌─────────▼─────────┐
              │   INTEGRATIONS    │
              │                   │
              │ Email             │
              │ Digital signature │
              │ Storage           │
              └───────────────────┘
```

**Estado de referencia:** rama `itinerary`, commit `d5339586cfe5fc01e9b218670a079bba6d2c82be`.