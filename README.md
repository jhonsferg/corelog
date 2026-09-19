# CoreLog
Un minisistema para la gestión de incidentes y la resolución de tickets.
Este proyecto se está desarrollando con fines educativos como parte del curso "Herramientas de Desarrollo" de la Universidad Tecnológica del Perú.

## Tabla de contenidos

- [Arquitectura](#architecture)
- [Tecnologías del backend](#backend-stack)
- [Tecnologías del frontend](#frontend-stack)
- [Requisitos previos](#prerequisites)
- [Primeros pasos](#getting-started)
- [Scripts de configuración](#setup-scripts)
- [Variables de entorno](#environment-variables)
- [Referencia de la API](#api-reference)
- [Comandos disponibles](#available-commands)
- [Docker](#docker)
- [Convenciones del proyecto](#project-conventions)
- [Solución de problemas](#troubleshooting)

## Arquitectura

Tanto el backend como el frontend siguen una arquitectura hexagonal, es decir, basada en puertos y adaptadores. La lógica de negocio, incluido el dominio y los casos de uso, no tiene conocimiento de mecanismos de entrega como HTTP, SQL, React o Axios. Los adaptadores se conectan al núcleo mediante interfaces explícitas llamadas puertos, lo que mantiene el sistema testeable y permite intercambiar cualquier adaptador, por ejemplo PostgreSQL por SQLite o REST por GraphQL, sin modificar las reglas de negocio.


```

corelog/
├── cmd/api/                    Punto de composición de Go, conecta todo e inicia el servidor HTTP
├── internal/
│   ├── platform/                infraestructura transversal: configuración, logger, controladores de BD, servidor HTTP, JWT, middleware
│   ├── user/                    hexágono de usuario para autenticación e identidad
│   │   ├── domain/                 entidades y reglas de negocio, sin dependencias externas
│   │   ├── application/
│   │   │   ├── port/                 in contiene interfaces de casos de uso, out contiene interfaces de repositorio
│   │   │   └── usecase/              implementaciones de casos de uso, el núcleo de la aplicación
│   │   ├── adapter/
│   │   │   ├── in/http/              adaptador de entrada: handlers HTTP, DTOs y rutas
│   │   │   ├── out/postgres/         adaptador de salida: repositorio basado en sqlx para producción
│   │   │   └── out/sqlite/           adaptador de salida: repositorio basado en sqlx para desarrollo local
│   │   └── module.go                 conecta el dominio, la aplicación y el adaptador
│   ├── team/                    hexágono de equipos para agrupar usuarios, misma estructura que user
│   └── ticket/                  hexágono de tickets para incidentes, misma estructura que user, además de adapter/out/memory y adapter/out/directory
├── migrations/
│   ├── postgres/                 migraciones SQL versionadas para producción, aplicadas con golang-migrate
│   └── sqlite/                   migraciones SQL versionadas para desarrollo local, integradas en el binario y aplicadas automáticamente al iniciar
├── docker/                      Dockerfiles y configuración de nginx
├── docker-compose.yml            PostgreSQL y API por defecto; agrega el perfil full para una UI en contenedor
├── openapi.yaml                  Contrato AbreAPI 3 para toda la API, mantenido sincronizado con la referencia de API de abajo
├── scripts/                      setup.sh and setup.ps1, comandos de configuración del entorno, compilación y empaquetado
└── ui/                           frontend React y TS, con estructura hexagonal, documentado en ui/README.md

```

### Flujo de solicitud, ejemplo: crear ticket

`HTTP request -> adapter/in/http as the Handler -> application/port/in as the Service interface -> application/usecase for business orchestration -> domain for Ticket invariants -> application/port/out as the Repository interface -> adapter/out/postgres or adapter/out/sqlite via sqlx`

## Tecnologías del backend

- **Lenguaje**: Go
- **Router**: [chi](https://github.com/go-chi/chi)
- **Configuración**: [godotenv](https://github.com/joho/godotenv) loads `.env` into the process environment, then [caarlos0/env](https://github.com/caarlos0/env) parses that environment straight into typed config structs using struct tags, so `internal/platform/config` has no manual parsing code
- **Acceso a la base de datos**: `database/sql` plus [sqlx](https://github.com/jmoiron/sqlx), no ORM, explicit SQL
- **Autenticación**: JWT via [golang-jwt](https://github.com/golang-jwt/jwt) plus bcrypt password hashing
- **Validación**: [go-playground/validator](https://github.com/go-playground/validator)
- **Migraciones**: SQL files under `migrations/postgres/`, run with [golang-migrate](https://github.com/golang-migrate/migrate)

### Dos bases de datos, un solo conjunto de reglas de negocio

`internal/user/module.go` and `internal/ticket/module.go` both take a `port.Repository` by injection instead of building one themselves. The composition root, `cmd/api/main.go`, decides which driven adapter to wire in, based on `DB_DRIVER`:

| `DB_DRIVER` | Uso previsto | Adapter package | Esquema |
| --- | --- | --- | --- |
| `postgres`, the default | Producción | `adapter/out/postgres` | Tracked migrations in `migrations/postgres/`, applied manually with `golang-migrate` |
| `sqlite` | Desarrollo local | `adapter/out/sqlite` | Tracked migrations in `migrations/sqlite/`, embedded into the binary via `go:embed` and applied automatically on every startup, zero setup, no server required |

The SQLite driver is [`modernc.org/sqlite`](https://modernc.org/sqlite), a pure-Go implementation with no CGO requirement, so `go run ./cmd/api` works out of the box on any platform without a C toolchain. Both adapters implement the exact same `port.Repository` interface per hexagon, so the domain, use cases, and HTTP handlers never change based on which database is active. The SQL itself differs between `migrations/postgres/` and `migrations/sqlite/` because the two engines use different column types, for example `UUID` and `TIMESTAMPTZ` in PostgreSQL against `TEXT` in SQLite, but each numbered pair of files creates the same two tables with the same constraints.

The ticket hexagon additionally supports `TICKET_REPOSITORY=memory`, which forces `adapter/out/memory`, an in-memory, pre-seeded, thread-safe repository, regardless of `DB_DRIVER`. This is useful for demos or tests that should not touch disk at all.

### Equipos y asignación

Every user optionally belongs to one team, tracked as a nullable `team_id` on the user. Teams themselves are a small, independent hexagon under `internal/team`, managed by administrators through the `/teams` endpoints.

When a ticket is assigned, `internal/ticket/application/usecase/ticket_service.go` checks that the user performing the assignment and the user being assigned share the same team, using an outbound port, `ticket/application/port.UserDirectory`, implemented once in `internal/ticket/adapter/out/directory` for both PostgreSQL and SQLite via `sqlx.DB.Rebind`. A user with no team can neither assign tickets nor be assigned one, and a mismatched team returns `403`. This keeps the ticket hexagon decoupled from the user hexagon's internal types: it only depends on a narrow interface it declares itself.

There is no dedicated admin bootstrap step: the very first account ever registered becomes an admin automatically, and every account after that defaults to requester. Admins manage teams and assign other users to them from the Teams screen in the frontend.

## Tecnologías del frontend

See [`ui/README.md`](ui/README.md) for the frontend's hexagonal layout. Stack: React plus TypeScript and Vite, React Router, Zustand, Axios, React Hook Form plus Zod, Ant Design, CSS Modules.

## Requisitos previos

- Go 1.27 or newer
- Node.js 24 or newer, and npm
- Docker and Docker Compose, needed only for the PostgreSQL production path
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI, needed only to run PostgreSQL migrations manually or to use the optional `make migrate-sqlite-*` commands; the application itself applies SQLite migrations automatically without this CLI

## Primeros pasos

Ejecutar CoreLog localmente significa ejecutar dos procesos al mismo tiempo: la API de Go y el frontend de React. Usa dos ventanas o pestañas de terminal separadas, una para el backend y otra para el frontend, y mantén ambos procesos ejecutándose mientras trabajas.

### 1. Iniciar el backend

Elige una de las dos opciones de base de datos siguientes y luego inicia la API. Ejecuta estos comandos desde la raíz del repositorio.

#### Option A: SQLite, configuración cero, recomendado para desarrollo local

```sh
cp .env.example .env

```

`.env.example` already sets `DB_DRIVER=sqlite`, so no further edits to `.env` are required. Inicia la API:

```sh
make run

```

La API crea `./data/<DB_NAME>.db`, `./data/corelog.db` with the default `.env.example` value, and its schema automatically on first run. No Docker, no Postgres, and no migration step are required. Deja este proceso ejecutándose: la API ahora está escuchando en `http://localhost:8080`.

#### Option B: PostgreSQL mediante Docker, similar a producción

```sh
cp .env.example .env

```

Edit `.env` and set `DB_DRIVER=postgres`, luego inicia los contenedores de la base de datos y de la API y aplica las migraciones versionadas:

```sh
make docker-up
make migrate-up

```

La API ahora está escuchando en `http://localhost:8080`.

### 2. Iniciar el frontend

Abre una segunda terminal, mantén el backend del paso 1 ejecutándose y, desde la raíz del repositorio, ejecuta:

```sh
cd ui
cp .env.example .env
npm install
npm run dev

```

`ui/.env.example` sets `VITE_API_BASE_URL=http://localhost:8080/api/v1`, matching the backend from step 1, so no further edits are required for a default local setup. Deja este proceso ejecutándose: el servidor de desarrollo del frontend ahora está escuchando en `http://localhost:5173`.

### 3. Verificar que ambos estén ejecutándose

Abre `http://localhost:5173` in a browser, registra una nueva cuenta, inicia sesión y crea un ticket. La página se comunica con la API at `http://localhost:8080` in the background.

Para comprobar el backend por separado, sin el navegador:

```sh
curl http://localhost:8080/health

```

Esto devuelve `{"status":"ok"}` cuando la API es accesible.

## Scripts de configuración

`scripts/setup.sh`, for bash, and `scripts/setup.ps1`, for PowerShell, wrap the steps above into single commands. Both scripts expose the same commands and behave the same way, so pick whichever matches your shell.

| Comando | Descripción |
| --- | --- |
| `check` | Indica si Go, Node, npm, Docker y golang-migrate están disponibles |
| `env` | Create `.env` and `ui/.env` from their example files if missing, without overwriting existing ones |
| `dev` | Prepare a local development environment backed by SQLite: runs `env`, downloads Go modules, installs frontend dependencies. No Docker required |
| `prod` | Prepare a similar a producción environment backed by PostgreSQL via Docker: verifies Docker is available, runs `env`, sets `DB_DRIVER=postgres` in a freshly created `.env`, starts the `db` and `api` containers, and applies migrations when the `migrate` CLI is installed |
| `build` | Compila el binario del backend en `bin/api`, `bin/api.exe` on Windows, y compila el frontend en `ui/dist` |
| `package` | Run `build`, then assemble a distributable archive under `dist/` que contiene el binario del backend, las migraciones de PostgreSQL y la compilación del frontend, y comprímelo como `corelog-package.tar.gz` in bash or `corelog-package.zip` in PowerShell |
| `clean` | Elimina `bin/`, `ui/dist/`, `dist/`, and the packaged archive |

Ejemplos:

```sh
./scripts/setup.sh dev
./scripts/setup.sh prod
./scripts/setup.sh package

```

```powershell
.\scripts\setup.ps1 dev
.\scripts\setup.ps1 prod
.\scripts\setup.ps1 package

```

Ejecutar cualquiera de los scripts sin un comando o con uno no reconocido muestra la lista de comandos anterior.

## Variables de entorno

All variables live in `.env`, copied from `.env.example`. See [`ui/.env.example`](https://www.google.com/search?q=ui/.env.example&utm_source=gemini) para ver las variables del frontend.

| Variable | Predeterminado | Descripción |
| --- | --- | --- |
| `APP_ENV` | `development` | `development` enables verbose text logs; anything else logs JSON |
| `SERVER_PORT` | `8080` | Puerto HTTP en el que escucha la API |
| `SERVER_READ_TIMEOUT` and `SERVER_WRITE_TIMEOUT` | `10s` | Tiempos de espera del servidor HTTP |
| `SERVER_SHUTDOWN_TIMEOUT` | `15s` | Período de gracia para solicitudes en curso durante el apagado |
| `CORS_ALLOWED_ORIGIN` | `http://localhost:5173` | Lista separada por comas de orígenes autorizados para llamar a la API desde el navegador |
| `DB_DRIVER` | `sqlite` in the example file, `postgres` as the code default | `postgres` or `sqlite` |
| `DB_NAME` | `corelog` | Nombre de base de datos independiente del motor. Para PostgreSQL, es la base de datos dentro del servidor. Para SQLite, es el archivo `data/<DB_NAME>.db`; there is no separate SQLite-only path variable |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE` | `localhost`, `5432`, `corelog`, `corelog`, `disable` | Conexión a PostgreSQL, utilizada solo cuando `DB_DRIVER=postgres` |
| `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME` | `25`, `25`, `5m` | Configuración del pool de conexiones de PostgreSQL |
| `JWT_SECRET` | `change-me-in-production` | Secreto de firma HMAC para los tokens de acceso; establece un secreto real fuera del desarrollo local |
| `JWT_EXPIRATION` | `24h` | Duración del token de acceso |
| `TICKET_REPOSITORY` | `postgres` | `postgres` follows `DB_DRIVER` despite the name, or `memory` forces the in-memory ticket adapter |

## Referencia de la API

El contrato completo se encuentra en [`openapi.yaml`](https://www.google.com/search?q=openapi.yaml&utm_source=gemini) at the repository root. La tabla siguiente es un resumen.

Ruta base: `/api/v1`. Los endpoints que requieren autenticación esperan `Authorization: Bearer <token>`.

| Método | Path | Autenticación | Descripción |
| --- | --- | --- | --- |
| `POST` | `/auth/register` | No requerida | Crea una cuenta de usuario. La primera cuenta creada se convierte automáticamente en administrador; todas las cuentas posteriores son solicitantes por defecto |
| `POST` | `/auth/login` | No requerida | Autentica al usuario y recibe un JWT |
| `GET` | `/users/me` | Requerida | Usuario autenticado actual |
| `PATCH` | `/users/me` | Requerida | Actualiza el nombre y correo electrónico del usuario actual |
| `POST` | `/users/me/password` | Requerida | Cambia la contraseña del usuario actual proporcionando la contraseña vigente |
| `GET` | `/users/search` | Requerida | Busca usuarios dentro del propio equipo del solicitante mediante `?q=`, for picking a ticket assignee |
| `GET` | `/users` | Solo administrador | Lista o busca todos los usuarios mediante `?q=` and `?team_id=` |
| `PATCH` | `/users/{id}/team` | Solo administrador | Asigna un usuario a un equipo o lo desvincula con `team_id: null` |
| `GET` | `/teams` | Requerida | Lista todos los equipos |
| `POST` | `/teams` | Solo administrador | Crea un equipo |
| `GET` | `/teams/{id}` | Requerida | Obtiene un equipo por ID |
| `PATCH` | `/teams/{id}` | Solo administrador | Cambia el nombre de un equipo |
| `DELETE` | `/teams/{id}` | Solo administrador | Elimina un equipo; sus miembros quedan sin equipo en lugar de ser eliminados |
| `POST` | `/tickets` | Requerida | Crea un ticket |
| `GET` | `/tickets` | Requerida | Lista tickets, filtrables mediante `?status=`, `?priority=`, `?page=`, `?page_size=` |
| `GET` | `/tickets/{id}` | Requerida | Obtiene un ticket por ID |
| `PATCH` | `/tickets/{id}/status` | Requerida | Cambia el estado de un ticket: `open → in_progress → resolved → closed`, with `resolved → in_progress` also allowed |
| `PATCH` | `/tickets/{id}/assign` | Requerida | Asigna un ticket a un usuario; el asignado debe pertenecer al mismo equipo que el solicitante; de lo contrario, devuelve `403` |
| `GET` | `/health` | No requerida | Liveness check, mounted at the root path without the `/api/v1` prefix |

Ticket `status`: `open`, `in_progress`, `resolved`, `closed`. Ticket `priority`: `low`, `medium`, `high`, `critical`. Admin-only endpoints require the caller's JWT to carry the `admin` role.

## Comandos disponibles

Backend, from the `Makefile` at the repo root:

| Comando | Descripción |
| --- | --- |
| `make run` | Ejecuta la API localmente mediante `go run ./cmd/api` |
| `make build` | Compila el binario de la API en `bin/api` |
| `make test` | Ejecuta las pruebas de Go |
| `make lint` | Run `go vet` |
| `make tidy` | Run `go mod tidy` |
| `make docker-up`, `make docker-down`, `make docker-logs` | Administra el stack de Docker Compose de PostgreSQL y la API |
| `make migrate-up`, `make migrate-down` | Aplica o revierte las migraciones de PostgreSQL |
| `make migrate-create name=<name>` | Genera el par de archivos de una nueva migración de PostgreSQL |
| `make migrate-sqlite-up`, `make migrate-sqlite-down` | Aplica o revierte las migraciones de SQLite sobre `data/<DB_NAME>.db` using the `migrate` CLI directly, as an alternative to the automatic startup migration; pass `DB_NAME=<name>` to target a database other than the default `corelog` |
| `make migrate-sqlite-create name=<name>` | Genera el par de archivos de una nueva migración de SQLite |
| `make ui-install`, `make ui-dev`, `make ui-build`, `make ui-lint` | Reenvía a los scripts npm del frontend |

Frontend, from `ui/package.json`, run inside `ui/`: `npm run dev`, `npm run build`, `npm run lint`, `npm run preview`.

## Docker

`docker-compose.yml` always starts `db`, which is PostgreSQL, and `api`. The frontend is optional and only started with the `full` profile:

```sh
docker compose up -d --build
docker compose --profile full up -d --build

```

The first command starts `db` and `api`. The second additionally starts `ui`, an nginx container serving the static production build.

### Contenedor de base de datos: usuario, base de datos y permisos

The `db` service runs the official `postgres:17-alpine` image unmodified. That image creates the role, the database, and the ownership grant between them automatically, driven entirely by three environment variables set in `docker-compose.yml`:

| Variable en el contenedor | Sourced from `.env` | Efecto |
| --- | --- | --- |
| `POSTGRES_USER` | `DB_USER`, default `corelog` | Nombre del rol de PostgreSQL que la imagen crea en el primer inicio. A este rol se le concede el privilegio de superusuario dentro del contenedor. |
| `POSTGRES_PASSWORD` | `DB_PASSWORD`, default `corelog` | Contraseña establecida para ese rol. |
| `POSTGRES_DB` | `DB_NAME`, default `corelog` | Name of the database the image creates, owned by `POSTGRES_USER`. |

Como el rol es propietario de la base de datos que crea, ya tiene permisos completos de lectura y escritura. No se requiere manualmente `CREATE ROLE`, `CREATE DATABASE`, or `GRANT` para la configuración predeterminada descrita en este proyecto.

This automatic setup runs only once: the first time the container starts against an empty `corelog_db_data` volume. Changing `DB_USER`, `DB_PASSWORD`, or `DB_NAME` in `.env` after that volume already holds data has no effect on the existing role or database, since the image only runs its initialization scripts against an empty data directory. To apply new credentials from a clean state, remove the volume before starting the containers again:

```sh
docker compose down
docker volume rm corelog_corelog_db_data
docker compose up -d --build

```

El nombre del volumen es el nombre del proyecto de Compose, normalmente el nombre del directorio del repositorio `corelog`, seguido de un guion bajo y del nombre del volumen declarado en `docker-compose.yml`, `corelog_db_data`. Run `docker volume ls` first to confirm the exact name on your machine, since a different directory name changes the project name prefix.

### Conectarse manualmente a la base de datos

Abre an interactive `psql` shell inside the running container:

```sh
docker compose exec db psql -U corelog -d corelog

```

Reemplaza `corelog` en ambos lugares por tus propios `DB_USER` and `DB_NAME` si los cambiaste respecto de los valores predeterminados. Desde fuera del contenedor, usando una instalación local de `psql` installation o un cliente gráfico como pgAdmin o TablePlus, conéctate con:

```sh
psql "host=localhost port=5432 user=corelog password=corelog dbname=corelog sslmode=disable"

```

Once connected, these `psql` meta-commands are useful for inspecting the setup:

```sql
\du
\l
\dt

```

`\du` lista los roles y sus atributos, `\l` lista las bases de datos y `\dt` lista las tablas de la base de datos actual.

### Crear manualmente un rol adicional

La configuración predeterminada proporciona a la aplicación un único rol con propiedad total de su base de datos, lo cual es suficiente para el desarrollo local y el alcance de este proyecto. Para crear un segundo rol, por ejemplo una cuenta de reportes de solo lectura, conéctate como se describió anteriormente y ejecuta:

```sql
CREATE ROLE reporting WITH LOGIN PASSWORD 'change-me';
GRANT CONNECT ON DATABASE corelog TO reporting;
GRANT USAGE ON SCHEMA public TO reporting;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO reporting;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO reporting;

```

La última instrucción garantiza que las tablas creadas por futuras migraciones también sean legibles por `reporting`, y no solo las tablas que ya existían cuando se ejecutó el permiso.

SQLite no tiene un modelo equivalente de usuarios o permisos: `migrations/sqlite/` produces a single file at `data/<DB_NAME>.db`, y el acceso a ella está controlado completamente por los permisos de archivos del sistema operativo en esa ruta.

## Convenciones del proyecto

* No comments in source code. Names, small functions, and the hexagonal layering carry the intent; this file and `ui/README.md` are where narrative documentation belongs.
* Ports before adapters. Every use case depends on an interface declared in `application/port`, never directly on `adapter/out/*`. The composition root is the only place concrete adapters are chosen: `cmd/api/main.go` on the backend, `adapters/in/state/*.store.ts` on the frontend.
* Un adaptador, una responsabilidad de persistencia. `adapter/out/postgres` and `adapter/out/sqlite` cada uno se encarga de su propio mapeo de filas al dominio; nunca comparten una estructura, por lo que cualquiera puede cambiar su representación de almacenamiento de forma independiente.
* Frontend imports use the `@/` alias, mapped to `ui/src`, for anything crossing a top-level boundary such as `app`, `core`, `adapters`, or `shared`. Same-directory and one-level-up imports stay relative. See `ui/README.md`.

## Solución de problemas

* `database is locked` on SQLite: only one process should hold `./data/<DB_NAME>.db` a la vez. Detén cualquier otra instancia de la API que esté ejecutándose.
* `migrate: command not found`: install [golang-migrate](https://www.google.com/search?q=https://github.com/golang-migrate/migrate%2523installation&utm_source=gemini). Solo es necesario para la ruta de PostgreSQL, no para SQLite.
* Errores de CORS en el navegador: asegúrate de que `CORS_ALLOWED_ORIGIN` coincida con el origen exacto en el que se ejecuta el frontend, cuyo valor predeterminado es `http://localhost:5173`.
* `401 Unauthorized` on `/tickets`: the ticket and `/users/me` endpoints require `Authorization: Bearer <token>` from `POST /auth/login`.