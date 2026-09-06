# Capacidades del backend InariOps

Este documento describe exclusivamente la funcionalidad implementada en Go. No incluye HTML, JavaScript ni funcionalidades del frontend.

## 1. Información general

- API HTTP bajo el prefijo `/api/v1`.
- Endpoint de salud: `GET /health`.
- Servidor HTTP en el puerto `9142`.
- CORS habilitado.
- Soporte para peticiones `OPTIONS`.
- Respuestas JSON y handlers generales para errores `404` y `405`.
- Los endpoints no tienen actualmente un middleware global que valide JWT o permisos.

## 2. Autenticación

Rutas definidas en `internal/api/auth_routes.go`.

| Método | Ruta | Funcionalidad |
|---|---|---|
| `POST` | `/api/v1/auth/` | Iniciar sesión con email y contraseña. |
| `PATCH` | `/api/v1/auth/` | Cambiar la contraseña de un usuario. |

El login:

- Comprueba que el usuario exista y esté activo.
- Valida la contraseña mediante bcrypt.
- Genera un JWT.
- Devuelve el ID, nombre, email, rol y si debe cambiar la contraseña.

El cambio de contraseña:

- Comprueba la contraseña anterior.
- Exige una nueva contraseña con reglas de complejidad.
- Actualiza la contraseña almacenada.

No existe endpoint de logout ni recuperación de contraseña.

## 3. Usuarios

Rutas definidas en `internal/api/user_routes.go`.

| Método | Ruta | Funcionalidad |
|---|---|---|
| `GET` | `/api/v1/users` | Listar usuarios. |
| `POST` | `/api/v1/users` | Crear usuario. |
| `GET` | `/api/v1/users/guides` | Listar usuarios que son guías. |
| `GET` | `/api/v1/users/{id}` | Obtener un usuario por ID. |
| `PATCH` | `/api/v1/users/{id}` | Actualizar un usuario. |
| `DELETE` | `/api/v1/users/{id}` | Desactivar un usuario. |

La creación de usuarios:

- Acepta los roles `ADMIN` y `GUIDE`.
- Crea las credenciales del usuario.
- Para un guía, crea automáticamente su registro de guía.
- Inicializa `MaxToursPerDay` en `1`.

El borrado implementado para usuarios es lógico mediante desactivación. Existe código interno para borrar físicamente, pero no una ruta pública para ello.

## 4. Reservas

Rutas definidas en `internal/api/reservations_routes.go`.

| Método | Ruta | Funcionalidad |
|---|---|---|
| `GET` | `/api/v1/reservations` | Listar reservas. |
| `POST` | `/api/v1/reservations` | Crear una reserva completa. |
| `GET` | `/api/v1/reservations/{code}` | Obtener una reserva por código. |
| `GET` | `/api/v1/reservations/details/{code}` | Obtener el detalle completo de una reserva. |
| `PATCH` | `/api/v1/reservations/{code}` | Actualizar una reserva. |
| `PATCH` | `/api/v1/reservations/{code}/signature` | Actualizar el estado de firma de una reserva. |
| `PATCH` | `/api/v1/reservations/{code}/voucher` | Actualizar el estado de voucher de una reserva. |
| `DELETE` | `/api/v1/reservations/{code}` | Cancelar una reserva. |
| `DELETE` | `/api/v1/reservations/erase/{code}` | Borrar físicamente una reserva. |
| `GET` | `/api/v1/reservations/{code}/customers` | Listar clientes de una reserva. |

La creación de una reserva:

- Valida que la agencia exista.
- Genera el identificador y el código de la reserva.
- Inicializa el estado como `PENDING_ASSIGNMENT`.
- Inicializa la firma como `NOT_SENT`.
- Inicializa el voucher como `NOT_GENERATED`.
- Crea automáticamente un día de tour por cada fecha del intervalo de la reserva.
- Crea clientes nuevos o reutiliza clientes existentes por email.
- Crea las relaciones entre la reserva y sus clientes.

La actualización de una reserva:

- Actualiza campos parciales.
- Valida la agencia si se modifica.
- Sincroniza los clientes asociados.
- Añade clientes nuevos.
- Mantiene clientes existentes.
- Elimina las relaciones que ya no aparecen en la lista recibida.
- Si recibe una lista vacía de clientes, elimina todas las asociaciones.

La cancelación normal es lógica: cambia la reserva a `CANCELLED` y cancela sus días. La operación `erase` elimina físicamente la reserva, sus días y sus relaciones.

La firma y el voucher ya disponen de endpoints para actualizar sus estados. La generación del PDF, el envío de emails y la integración externa de firma siguen pendientes.

## 5. Días de tour

Rutas definidas en `internal/api/tourday_routes.go`.

| Método | Ruta | Funcionalidad |
|---|---|---|
| `PATCH` | `/api/v1/tour-days/assign-guide` | Asignar guía a uno o varios días. |
| `GET` | `/api/v1/tour-days/by-reservation/{reservation_id}` | Listar días de una reserva. |
| `GET` | `/api/v1/tour-days/available-for-guide/{guide_id}` | Listar días elegibles para un guía según disponibilidad y límite diario. |
| `POST` | `/api/v1/tour-days` | Crear un día de tour. |
| `GET` | `/api/v1/tour-days/{id}` | Obtener un día por ID. |
| `PATCH` | `/api/v1/tour-days/{id}` | Actualizar un día. |
| `PATCH` | `/api/v1/tour-days/{id}/confirm` | Confirmar explícitamente un día preasignado. |
| `DELETE` | `/api/v1/tour-days/{id}` | Cancelar un día. |

Un día de tour puede contener:

- Título.
- Fecha y hora.
- Duración.
- Zona.
- Guía.
- Número de personas.
- Punto de encuentro.
- Hotel del guía.
- Hotel del cliente.
- Notas de responsabilidad para guía y cliente.
- Remuneración.
- Estado.
- Estado del voucher.

La creación inicializa el día como `PENDING_ASSIGNMENT` y el voucher como `NOT_GENERATED`.

La asignación de guía:

- Puede aplicarse a uno o varios días.
- Cambia el estado a `GUIDE_PREASSIGNED`.
- Registra historial de asignación.

Existe un servicio y una ruta separada para desasignar guías y devolver el día a `PENDING_ASSIGNMENT`.

La cancelación es lógica y cambia el día a `CANCELLED`.

## 6. Guías

Rutas definidas en `internal/api/guides_routes.go`.

| Método | Ruta | Funcionalidad |
|---|---|---|
| `GET` | `/api/v1/guides` | Listar guías con información enriquecida. |
| `POST` | `/api/v1/guides` | Crear una guía completa. |
| `GET` | `/api/v1/guides/{id}` | Obtener una guía por ID. |
| `PATCH` | `/api/v1/guides/{id}` | Actualizar una guía. |
| `GET` | `/api/v1/guides/user/{user_id}` | Obtener una guía por usuario. |
| `GET` | `/api/v1/guides/{id}/languages` | Listar idiomas de una guía. |
| `POST` | `/api/v1/guides/{id}/languages` | Añadir un idioma a una guía. |
| `DELETE` | `/api/v1/guides/{id}/languages/{language_id}` | Eliminar un idioma de una guía. |

La creación completa:

- Crea un usuario con rol `GUIDE`.
- Crea el registro de guía.
- Asocia idiomas.
- Asocia zonas.
- Devuelve información enriquecida.

La actualización puede sincronizar completamente los idiomas y las zonas: elimina asociaciones ausentes y añade las nuevas.

Existe funcionalidad interna para añadir zonas individualmente, pero no hay un endpoint HTTP específico para ello.

## 7. Clientes

Rutas definidas en `internal/api/customer_routes.go`.

| Método | Ruta | Funcionalidad |
|---|---|---|
| `GET` | `/api/v1/customers` | Listar clientes. |
| `POST` | `/api/v1/customers` | Crear cliente. |
| `GET` | `/api/v1/customers/{id}` | Obtener cliente por ID. |
| `PATCH` | `/api/v1/customers/{id}` | Actualizar cliente. |
| `DELETE` | `/api/v1/customers/{id}` | Eliminar cliente. |

El sistema también puede:

- Buscar clientes por nombre.
- Buscar por documento.
- Buscar por teléfono.
- Buscar por email.
- Limitar la búsqueda a 20 resultados.
- Crear o reutilizar clientes automáticamente al crear reservas.

La búsqueda avanzada existe internamente, pero no está conectada a una ruta HTTP propia.

## 8. Agencias

Rutas definidas en `internal/api/agency_routes.go`.

| Método | Ruta | Funcionalidad |
|---|---|---|
| `GET` | `/api/v1/agencies` | Listar agencias. |
| `POST` | `/api/v1/agencies` | Crear agencia. |
| `GET` | `/api/v1/agencies/{id}` | Obtener agencia por ID. |
| `PATCH` | `/api/v1/agencies/{id}` | Actualizar agencia. |
| `DELETE` | `/api/v1/agencies/{id}` | Eliminar agencia. |

La lista admite filtros por:

- `name`.
- `email`.
- `is_active`.

La creación:

- Genera el ID.
- Comprueba que no exista el email.
- Fuerza la agencia a estado activo inicialmente.

## 9. Idiomas

Rutas definidas en `internal/api/language_routes.go`.

| Método | Ruta | Funcionalidad |
|---|---|---|
| `GET` | `/api/v1/languages` | Listar idiomas. |
| `POST` | `/api/v1/languages` | Crear idioma. |
| `GET` | `/api/v1/languages/{code}` | Obtener idioma por código. |
| `PATCH` | `/api/v1/languages/{code}` | Actualizar idioma. |
| `DELETE` | `/api/v1/languages/{code}` | Eliminar idioma. |

Los idiomas nuevos se crean activos y no se permiten códigos duplicados.

## 10. Zonas

Rutas definidas en `internal/api/zone_routes.go`.

| Método | Ruta | Funcionalidad |
|---|---|---|
| `GET` | `/api/v1/zones` | Listar zonas. |
| `POST` | `/api/v1/zones` | Crear zona. |
| `GET` | `/api/v1/zones/{name}` | Obtener zona por nombre. |
| `PATCH` | `/api/v1/zones/{id}` | Actualizar zona. |
| `DELETE` | `/api/v1/zones/{name}` | Eliminar zona. |

Las zonas nuevas se crean activas y no se permiten nombres duplicados.

## 11. Estados del dominio

### Estados de reserva y día

- `PENDING`
- `PENDING_ASSIGNMENT`
- `GUIDE_PREASSIGNED`
- `PAYMENT_PENDING`
- `PARTIALLY_CONFIRMED`
- `CONFIRMED`
- `COMPLETED`
- `CANCELLED`
- `FORCE_MAJEURE_CANCELLED`

### Estados de firma

- `NOT_SENT`
- `SENT`
- `SIGNED`
- `REJECTED`

### Estados de voucher

- `NOT_GENERATED`
- `GENERATED`
- `PARTIALLY_SENT`
- `SENT`

### Transiciones implementadas

- Crear reserva o día: `PENDING_ASSIGNMENT`.
- Asignar guía: `GUIDE_PREASSIGNED`.
- Confirmar día preasignado: `GUIDE_CONFIRMED`.
- Cuando todos los días están `GUIDE_CONFIRMED`, sincronizar la reserva a `CONFIRMED`.
- Desasignar guía: `PENDING_ASSIGNMENT`.
- Cancelar reserva o día: `CANCELLED`.
- Expirar una asignación: `PAYMENT_PENDING`.
- Sincronizar el estado de la reserva a partir de sus días.

Existe un resolver de estados capaz de devolver un estado común o `PENDING` cuando los días tienen estados diferentes, aunque no está completamente conectado al flujo actual.

## 12. Workers automáticos

Los workers se inician automáticamente con el servidor y se ejecutan aproximadamente cada minuto.

### Expiración de asignaciones

Busca días en estado `GUIDE_PREASSIGNED` cuya actualización tenga más de 24 horas.

Cuando los encuentra:

- Los cambia a `PAYMENT_PENDING`.
- Registra el historial.
- Usa `SYSTEM` como actor.
- Guarda la razón `Guide assignment expired`.

### Sincronización de reservas

- Carga las reservas.
- Obtiene los días de cada reserva.
- Recalcula el estado general.
- Actualiza la reserva si el estado cambió.

## 13. Funcionalidad interna sin endpoint público

Existe código para las siguientes operaciones, pero no hay una ruta HTTP pública conectada:

- Búsqueda avanzada de clientes.
- Borrado físico de usuarios.
- Obtener detalle enriquecido de una guía individual.
- Añadir una zona individual a una guía.
- Añadir clientes explícitamente a una reserva.
- Eliminar clientes explícitamente de una reserva.
- Registrar manualmente historial de días.
- Borrar físicamente un día de tour.
- Obtener reservas por ID de agencia; la implementación actual devuelve `ErrNotFound`.
- Logout.
- Recuperación o restablecimiento de contraseña.
- Integración externa completa de firma, incluyendo envío y callback.
- Generación de vouchers PDF desde plantilla HTML y envío por email.
- Paginación y filtros avanzados de reservas.

## 14. Incidencias conocidas de rutas

- Algunas operaciones de actualización usan el ID enviado en el cuerpo en lugar del parámetro de la URL.
- La configuración JWT utiliza un secreto definido directamente en el código.
- Aunque se generan JWT, no hay validación global del header `Authorization`.

## 15. Comparación con el alcance objetivo

La estimación se refiere únicamente al backend actual. No se han contado HTML, JavaScript ni funcionalidades simuladas del frontend.

### Método de estimación

Se evaluaron 27 capacidades funcionales atómicas:

- Implementada: flujo HTTP y persistencia funcional.
- Parcial: existe algún modelo, campo o parte del flujo, pero no la funcionalidad completa.
- Ausente: no existe implementación backend identificable.

Resultado:

- 6 capacidades implementadas.
- 8 capacidades parciales.
- 13 capacidades ausentes.

Contando una capacidad parcial como media capacidad:

$$
\\frac{6 + 8 \\times 0.5}{27} \\times 100 = 37\\%
$$

La cobertura funcional del backend Go sigue siendo de aproximadamente **37%**, con un rango razonable de **35% a 40%**. Esta cifra mide flujos disponibles mediante servicios y API, no la existencia de tablas.

El ERD mejora la situación de preparación técnica: una parte importante de las capacidades antes consideradas ausentes ya tiene tablas, relaciones y estados definidos. Por eso conviene distinguir:

- **Cobertura funcional backend:** aproximadamente 37%.
- **Cobertura de persistencia/modelo:** sensiblemente mayor, porque el ERD ya prepara disponibilidad, itinerarios, incidencias, archivos, firmas, vouchers, notificaciones y auditoría.

La diferencia restante está principalmente en implementar los servicios, repositorios, handlers, rutas, integraciones de email/PDF y reglas de negocio que conecten ese modelo con la aplicación.

### Roles del sistema

| Requisito | Estado | Observación |
|---|---|---|
| Roles `ADMIN` y `GUIDE` | Parcial | Los roles existen y se incluyen en el JWT, pero no se validan tokens ni permisos en las rutas. |
| Operación administrativa protegida | Ausente | No existe autorización por rol. |
| Operación propia del guía protegida | Ausente | No existen endpoints con control de identidad del guía autenticado. |

### Modelo principal

| Requisito | Estado | Observación |
|---|---|---|
| Reserva padre | Implementado | CRUD, detalle, cancelación y borrado físico. |
| Días de tour hijos | Implementado | Una reserva puede contener varios días y se crean automáticamente según el rango de fechas. |
| Clientes | Implementado | CRUD, asociación a reservas y reutilización por email. |
| Agencias | Implementado | CRUD y filtros básicos. |
| Guías | Implementado | Gestión de usuario-guía, idiomas y zonas. |
| Disponibilidad de guías | Implementado parcialmente | Existen rutas HTTP para crear, consultar, actualizar y eliminar disponibilidades, con validación de fechas, rangos y solapamientos. La regla de negocio ya se aplica en la consulta de tours elegibles: un guía no disponible no aparece como candidato. Falta restringir el borrado al guía de la URL. |
| Itinerarios | Parcial | Existen datos operativos del día, pero no una entidad de itinerario con actividades ordenadas. |

### Documentación y procesos legales

| Requisito | Estado | Observación |
|---|---|---|
| Firma legal | Parcial | El modelo de datos ya contempla `legal_signatures`, estados, enlace externo y fechas; falta el flujo backend para enviar automáticamente el email y actualizar el estado cuando corresponda. |
| Bono o voucher PDF | Parcial | El modelo de datos ya contempla `vouchers`, versión, URL y fecha de envío; falta generar el PDF desde una plantilla HTML, almacenarlo y enviarlo por email. |
| Documentación asociada al tour | Parcial | Existe una base genérica de archivos (`files`), pero todavía no hay flujo backend de documentos asociados, generación o descarga. |

### Operación e incidencias

| Requisito | Estado | Observación |
|---|---|---|
| Crear, modificar y cancelar tours | Implementado | Se cubre mediante reservas y días de tour. |
| Confirmar tours | Parcial | Existen estados de confirmación, pero no existe un endpoint específico ni una transición completa de confirmación. |
| Tours disponibles para guías | Parcial | Existe una consulta de días elegibles por guía que excluye sus períodos de indisponibilidad y respeta `MaxToursPerDay`; todavía falta completar el flujo de aceptación por parte del guía. |
| Aceptar un tour como guía | Ausente | La asignación actual la realiza un operador; no existe aceptación del guía. |
| Cancelar o rechazar un tour como guía | Ausente | No existe operación con reglas de cancelación del guía. |
| Calendario de guías | Parcial | Se guardan fechas de tours, pero no existe consulta de calendario por guía o rango. |
| Historial de tours del guía | Parcial | Hay historial técnico de asignaciones y estados, pero no consulta funcional para el guía. |
| Incidencias | Parcial | La base de datos ya contempla `incidents`, estados e historial; falta el módulo Go con endpoints, permisos, transiciones y listado de incidencias actuales. |
| Fotografías de incidencias | Parcial | Existen `files` e `incident_files`; falta la carga y asociación desde la API. |

### Plataforma móvil y comunicación

| Requisito | Estado | Observación |
|---|---|---|
| Perfil del guía | Implementado | Existe información de usuario, idiomas, zonas y límite diario. |
| Notificaciones push/email | Parcial | La base de datos ya contempla `notifications`, canal PUSH/EMAIL y lectura; falta el servicio de notificaciones, proveedor y endpoints o eventos que las generen. |
| Estadísticas administrativas | Ausente | No hay endpoints ni consultas de métricas. |
| Auditoría de negocio | Parcial | La base de datos ya contempla `audit_logs`; falta registrar las acciones desde los servicios y exponer consultas si son necesarias. |

## 16. Diferencias entre el modelo solicitado y el modelo actual

### Entidades preparadas en la base de datos pero no implementadas en Go

El ERD confirma que varias capacidades no deben considerarse ausentes del diseño de persistencia. Existen tablas, relaciones o tipos para:

- `Itinerario`.
- `Bono` como documento versionado (`vouchers`).
- `FirmaLegal` como entidad de integración (`legal_signatures`).
- `Incidencia`.
- `Disponibilidad`.
- `Notificacion`.
- `LogAuditoria` de negocio.
- Archivos e imágenes mediante `files` e `incident_files`.

Estas piezas están preparadas en PostgreSQL, pero todavía no están conectadas a handlers, servicios, repositorios ni rutas HTTP del backend analizado.

La disponibilidad es una excepción parcial: ya tiene DTO, repositorio, servicio, handlers y rutas públicas en `internal/api/guides_routes.go`. El bloque usa `*time.Time`, permite `reason` vacío o `null`, acepta el mismo día y rechaza rangos inválidos o solapados. El JSON debe respetar el formato que acepta `time.Time` (por ejemplo, RFC3339); una fecha simple `YYYY-MM-DD` no se deserializa automáticamente con este contrato.

### Campos o conceptos parciales

- La agencia sí tiene el campo `representative` en el ERD, aunque todavía no está integrado de forma homogénea en todos los DTOs y flujos del backend.
- El ERD ya contempla clientes propios del día mediante `tour_day_customers` y actividades ordenadas mediante `itinerary_items`; todavía faltan sus módulos Go, servicios y rutas HTTP.
- El guía tiene zonas, idiomas, `MaxToursPerDay` y disponibilidad por rangos de fechas. La disponibilidad filtra la consulta de tours elegibles para el guía y respeta su límite diario. No forma parte de la responsabilidad de `assign-guide`; queda pendiente verificar la pertenencia del registro al guía en el borrado.
- El voucher ya tiene tabla propia con reserva/día, estado, versión, URL y fecha de envío; todavía falta la generación y distribución real.
- La firma ya tiene tabla propia con reserva, enlace externo, estado y fechas; todavía falta automatizar el email y el procesamiento del resultado.
- Los estados en inglés son coherentes con una API y un dominio escalables; no constituyen una diferencia funcional negativa respecto a los nombres en español de la guía.
- La confirmación manual por parte del guía debe tratarse como una regla de negocio a concretar. El ERD contempla `GUIDE_CONFIRMED`, mientras que el flujo Go actual debe conectarse explícitamente a esa transición.
- El voucher debe generarse a partir de la información almacenada, usando una plantilla HTML, convertirse a PDF y enviarse por email. El envío manual opcional puede servir como revisión operativa del formato antes de automatizarlo.
- En incidencias, el listado operativo debe mostrar únicamente estados activos (`OPEN` e `IN_REVIEW`). Al pasar a `RESOLVED` o `CLOSED`, la incidencia debe desaparecer del listado actual, conservándose en base de datos para historial y auditoría.

## 17. Qué falta para completar el alcance

### Prioridad alta

1. Implementar middleware JWT y autorización por rol.
2. Verificar la pertenencia del registro al guía al modificar o eliminar disponibilidades.
3. Implementar flujo de tours disponibles, aceptación y rechazo por parte del guía.
4. Crear el módulo de incidencias con estados, fotos, permisos, historial y filtro de incidencias activas (`OPEN`/`IN_REVIEW`). Las incidencias `RESOLVED`/`CLOSED` deben quedar fuera del listado operativo, sin borrarse físicamente.
5. Implementar firma legal: generación/envío del email, enlace o formato externo y actualización del estado.
6. Implementar vouchers: plantilla HTML, conversión a PDF, almacenamiento, revisión manual opcional y envío por email.

### Prioridad media

1. Crear itinerarios y actividades ordenadas por día de tour.
2. Añadir calendarios y consultas por guía, rango de fechas y zona.
3. Añadir estadísticas administrativas.
4. Crear notificaciones persistentes y notificaciones push.
5. Implementar auditoría de negocio con usuario, entidad, acción, descripción y timestamp.

### Correcciones necesarias del núcleo existente

1. Completar la validación de transiciones y permisos de confirmación.
2. Completar la generación/envío de firma y voucher.

## 18. Actualización incremental del análisis

Desde la revisión técnica anterior se completó el bloque HTTP y de validación de disponibilidad de guías:

- DTO `Availability` y solicitud de creación.
- Repositorio para crear, consultar, actualizar y eliminar registros.
- Servicio, handlers y rutas HTTP para las operaciones de disponibilidad.
- Validación de fechas mediante `*time.Time`.
- Rechazo de guía inexistente, fechas inválidas, rangos invertidos y rangos solapados.
- Aceptación de rangos de un solo día, rangos consecutivos y `reason` vacío o `null`.
- Comprobación de `RowsAffected` al actualizar o eliminar.

La capacidad queda **parcial** porque todavía no se consulta al asignar un tour ni se comprueba el solapamiento con `tour_days` asignados. Además, quedan estos riesgos concretos:

- `DELETE /guides/{id}/availabilities/{availability_id}` usa el `availability_id`, pero el servicio no verifica que ese registro pertenezca al guía `{id}`. Un cliente que conozca otro UUID podría eliminar una disponibilidad ajena.
- El borrado de disponibilidad no reutiliza el mapeo de errores del resto del bloque; una disponibilidad inexistente puede terminar como `500` en lugar de `404`.
- `HasConflictingAvailability` se consulta en el servicio, pero dos creaciones simultáneas pueden superar ambas la comprobación antes de insertar. Hace falta una estrategia transaccional o una restricción de exclusión en PostgreSQL.
- Hay tests unitarios para el rango de fechas, pero todavía faltan tests de handler, repositorio, conflicto real y permisos de pertenencia al guía.

La regla funcional es: dado un `tour_day.start_datetime`, excluir de la consulta de tours elegibles cualquier guía cuya disponibilidad cubra esa fecha. La consulta ya aplica ese filtro y también limita los resultados por `MaxToursPerDay`. `assign-guide` no necesita repetir esta validación porque opera sobre una asignación ya seleccionada.

## 19. Revisión técnica profunda

Esta revisión analiza problemas de seguridad, consistencia, diseño y robustez del backend Go. `go vet ./...` no detecta errores estáticos, pero eso no cubre los problemas funcionales ni de arquitectura descritos aquí.

### Críticos

#### 1. La API no autentica ni autoriza las peticiones

Se generan JWT durante el login, pero no existe middleware que valide `Authorization`, expiración, usuario o rol. Por tanto, las operaciones de creación, modificación, cancelación y borrado quedan expuestas a cualquier cliente que conozca las rutas.

**Impacto:** acceso no autorizado y posibilidad de modificar datos operativos completos.

**Recomendación:** añadir middleware JWT global y middleware de permisos por rol. Las operaciones propias del guía deben obtener el usuario desde el token, no desde un ID enviado por el cliente.

#### 2. Credenciales sensibles embebidas en el código

La conexión PostgreSQL contiene usuario, contraseña, IP y `sslmode=disable` directamente en `internal/db/postgres.go`. El secreto de firma JWT también está hardcodeado en `internal/modules/auth/jwt.go`.

**Impacto:** exposición de la base de datos y falsificación de tokens si el repositorio o una imagen Docker se filtra. La conexión además no cifra el tráfico.

**Recomendación:** mover secretos a variables de entorno o un gestor de secretos, exigir TLS en producción, rotar las credenciales actuales y configurar tiempos de conexión.

### Altos

#### 3. Creación de reservas sin transacción

`CreateReservation` crea primero la reserva, después sus días, clientes y relaciones en operaciones separadas. Si falla cualquiera de los pasos posteriores, queda una reserva incompleta en la base de datos.

**Impacto:** reservas sin días, clientes huérfanos o relaciones incompletas.

**Recomendación:** ejecutar toda la creación dentro de una transacción PostgreSQL y hacer rollback ante cualquier error.

#### 4. Actualización de reservas parcialmente aplicada

`UpdateReservation` actualiza primero la reserva y después sincroniza clientes con varias operaciones independientes. Un error al añadir o quitar una relación deja la reserva actualizada pero con clientes antiguos o incompletos.

**Impacto:** el detalle de la reserva deja de representar una operación atómica.

**Recomendación:** agrupar actualización y sincronización de clientes en una única transacción.

#### 5. Cambiar las fechas de una reserva no sincroniza sus días

La actualización permite cambiar `start_date` y `end_date`, pero no crea, elimina ni reajusta los registros de `tour_days` asociados.

**Impacto:** el rango de la reserva padre puede no coincidir con sus días hijos; el worker y la operación diaria trabajarán sobre información contradictoria.

**Recomendación:** definir una política explícita para cambios de rango y aplicarla transaccionalmente, protegiendo días ya asignados o confirmados.

#### 6. El worker puede sobrescribir estados terminales

`SyncReservationStatus` solo considera principalmente `PENDING_ASSIGNMENT`, `GUIDE_PREASSIGNED` y `PAYMENT_PENDING`. Una reserva `CONFIRMED`, `COMPLETED` o `CANCELLED` puede ser recalculada a `PENDING_ASSIGNMENT` si sus días no coinciden con los estados contemplados. Una reserva sin días también puede terminar en `PENDING_ASSIGNMENT`.

**Impacto:** regresiones automáticas de estado y pérdida de una cancelación o confirmación.

**Recomendación:** definir una máquina de estados explícita, excluir estados terminales del worker y tratar una reserva sin días como error o estado controlado.

#### 7. Borrado físico sin transacción y con dependencias incompletas

`EraseReservation` elimina días, relaciones y reserva mediante llamadas separadas. El ERD contiene dependencias adicionales como historiales, vouchers, firmas, incidencias y archivos, y esas relaciones no están configuradas de forma general con `ON DELETE CASCADE`.

**Impacto:** el borrado puede fallar a mitad de camino o quedar bloqueado por claves foráneas; también puede dejar registros relacionados sin limpiar.

**Recomendación:** evitar el borrado físico en operaciones normales; si se mantiene, usar una transacción y definir explícitamente la política de cascada o archivado para cada entidad.

#### 8. Asignación múltiple de guías no atómica

`AssignGuideToMultipleTourDays` actualiza días uno por uno sin transacción. Además, `AssignGuide` no comprueba inmediatamente el error de esa operación antes de continuar con el registro del historial.

**Impacto:** parte de los días puede quedar asignada y parte no, o puede existir asignación sin historial correspondiente.

**Recomendación:** usar una transacción para cambios e historial, comprobar cada error inmediatamente y validar límites de disponibilidad del guía.

#### 9. Cambio de contraseña ligado a un ID enviado por el cliente

`PATCH /auth/` recibe `user_id` en el cuerpo y no lo obtiene de una identidad autenticada. Actualmente no hay middleware que impida invocar el endpoint sin sesión.

**Impacto:** el contrato permite intentar operaciones sobre cualquier usuario; la autorización depende únicamente de conocer la contraseña anterior.

**Recomendación:** obtener el usuario del JWT. Para administración, crear un endpoint separado y protegido para restablecimiento de credenciales.

### Medios

#### 10. El `PATCH` de reservas ignora el código de la URL

La ruta es `/reservations/{code}`, pero `UpdateReservation` utiliza `Code` del JSON. Esto permite que el recurso indicado en la URL no coincida con el que realmente se actualiza y hace que el endpoint dependa de un campo interno del body.

**Recomendación:** obtener siempre el código con `mux.Vars(r)` y eliminarlo del DTO de entrada o rechazar conflictos entre ambos valores.

#### 11. Actualización de días sin validación de dominio

El servicio permite cambiar directamente `status`, `voucher_status`, guía, fecha y otros campos sin validar transiciones, existencia de relaciones ni reglas de negocio. También existe un archivo de validación de días sin implementación.

**Impacto:** se pueden producir estados imposibles, asignar guías inexistentes o modificar días cancelados/confirmados.

**Recomendación:** centralizar transiciones válidas y validar fechas, duración, personas, remuneración, guía, zona, hoteles y pertenencia a la reserva.

#### 12. Creación de usuarios crea una guía para cualquier rol

`CreateUser` crea siempre un registro en `guides`, incluso cuando el rol es `ADMIN`.

**Impacto:** datos incorrectos y ambigüedad sobre qué usuarios pueden recibir asignaciones.

**Recomendación:** crear el registro en `guides` únicamente para `GUIDE` y hacerlo dentro de la misma transacción que usuario y credenciales.

#### 13. El flujo de desasignación exige un guía que ya no existe

El servicio de desasignación recibe `guide_id` y usa `GetGuideByID` para crear el historial, aunque el día puede tener el guía guardado y el cliente no debería tener que reenviarlo. Si llega vacío o no coincide, la desasignación puede fallar después de haber actualizado el día.

**Recomendación:** leer el guía anterior antes de desasignar, registrar el historial con ese valor y hacer ambas operaciones atómicas.

#### 14. Errores HTTP demasiado genéricos

Muchos handlers convierten cualquier error en `500` o `401`, incluso cuando corresponde `400`, `404` o `409`. También se devuelven mensajes distintos según el handler para casos equivalentes.

**Impacto:** clientes incapaces de reaccionar correctamente y dificultad para diagnosticar fallos.

**Recomendación:** mapear errores de dominio a códigos HTTP uniformes y usar un formato de error JSON común.

### Diseño y operación

#### 15. El modelo de aplicación no aprovecha todavía el ERD

El ERD ya contiene disponibilidad, itinerarios, incidencias, archivos, firmas, vouchers, notificaciones y auditoría. La disponibilidad ya está registrada parcialmente en el bootstrap, pero las demás áreas todavía no tienen módulos o servicios conectados. La base de datos sigue por delante de la aplicación.

**Riesgo:** lógica futura duplicada o implementada directamente en handlers si no se mantiene la separación repository/service/application.

#### 16. Workers sin límites operativos ni apagado HTTP real

Los workers se lanzan en goroutines y reciben contexto, pero el servidor no usa `Shutdown` con señales del sistema, no cierra la conexión de base de datos y el worker de sincronización no pasa el contexto a sus consultas principales. Tampoco hay métricas, backoff ni alertas de fallos.

**Impacto:** apagados incompletos, consultas bloqueadas y errores silenciosos en producción.

**Recomendación:** implementar apagado coordinado, timeouts por consulta, métricas, logging estructurado y política de reintentos.

#### 17. Ausencia de pruebas de integración para invariantes críticas

Ya existe un test unitario para la validación del rango de disponibilidad. Aun así, faltan pruebas para autorización, rollback de creación/actualización, cambios de fechas, transiciones de estado, asignación múltiple, cancelación, solapamientos contra PostgreSQL y borrado con todas las claves foráneas del ERD.

**Recomendación:** añadir pruebas de repositorio con PostgreSQL y pruebas HTTP de integración antes de ampliar los módulos funcionales.

### Orden recomendado de corrección

1. Autenticación/autorización y retirada de secretos del código.
2. Transacciones para reservas, clientes, días, asignaciones y borrados.
3. Corrección del `PATCH /reservations/{code}` y de las transiciones de estado.
4. Validación completa de días, guías, fechas y relaciones.
5. Apagado controlado, timeouts y observabilidad.
6. Tests de integración sobre PostgreSQL.

