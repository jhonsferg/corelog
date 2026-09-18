# CoreLog v2

A mini system for incident management and ticket resolution.
This project is being developed for educational purposes as part of the "Development Tools" course at the Technological University of Peru.

## Table of contents

- [Architecture](#architecture)
- [Backend stack](#backend-stack)
- [Frontend stack](#frontend-stack)
- [Prerequisites](#prerequisites)
- [Getting started](#getting-started)
- [Setup scripts](#setup-scripts)
- [Environment variables](#environment-variables)
- [API reference](#api-reference)
- [Available commands](#available-commands)
- [Docker](#docker)
- [Project conventions](#project-conventions)
- [Troubleshooting](#troubleshooting)

## Architecture

Both the backend and the frontend follow hexagonal architecture, meaning ports and adapters. Business logic, including the domain and use cases, has zero knowledge of delivery mechanisms such as HTTP, SQL, React, or Axios. Adapters plug into the core through explicit interfaces called ports, which keeps the system testable and lets any adapter be swapped, for example PostgreSQL for SQLite or REST for GraphQL, without touching business rules.

```
corelog/
├── cmd/api/                    Go composition root, wires everything and starts the HTTP server
├── internal/
│   ├── platform/                cross-cutting infra: config, logger, db drivers, http server, jwt, middleware
│   ├── user/                    user hexagon for authentication and identity
│   │   ├── domain/                 entities and business rules, no external dependencies
│   │   ├── application/
│   │   │   ├── port/                 in holds use-case interfaces, out holds repository interfaces
│   │   │   └── usecase/              use-case implementations, the application core
│   │   ├── adapter/
│   │   │   ├── in/http/              driving adapter: HTTP handlers, DTOs, and routes
│   │   │   ├── out/postgres/         driven adapter: sqlx-based repository for production
│   │   │   └── out/sqlite/           driven adapter: sqlx-based repository for local development
│   │   └── module.go                 wires domain, application, and adapter together
│   ├── team/                    team hexagon for grouping users, same layout as user
│   └── ticket/                  ticket hexagon for incidents, same layout as user, plus adapter/out/memory and adapter/out/directory
├── migrations/
│   ├── postgres/                 tracked SQL migrations for production, applied with golang-migrate
│   └── sqlite/                   tracked SQL migrations for local development, embedded in the binary and auto-applied on startup
├── docker/                      Dockerfiles and nginx config
├── docker-compose.yml            PostgreSQL and API by default, add the full profile for a containerized UI
├── openapi.yaml                  OpenAPI 3 contract for the whole API, kept in sync with the API reference below
├── scripts/                      setup.sh and setup.ps1, environment setup, build, and packaging commands
└── ui/                           React and TS frontend, hexagonal layout, documented in ui/README.md
```

### Request flow, example: create ticket

`HTTP request -> adapter/in/http as the Handler -> application/port/in as the Service interface -> application/usecase for business orchestration -> domain for Ticket invariants -> application/port/out as the Repository interface -> adapter/out/postgres or adapter/out/sqlite via sqlx`

## Backend stack

- **Language**: Go
- **Router**: [chi](https://github.com/go-chi/chi)
- **Configuration**: [godotenv](https://github.com/joho/godotenv) loads `.env` into the process environment, then [caarlos0/env](https://github.com/caarlos0/env) parses that environment straight into typed config structs using struct tags, so `internal/platform/config` has no manual parsing code
- **Database access**: `database/sql` plus [sqlx](https://github.com/jmoiron/sqlx), no ORM, explicit SQL
- **Auth**: JWT via [golang-jwt](https://github.com/golang-jwt/jwt) plus bcrypt password hashing
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **Migrations**: SQL files under `migrations/postgres/`, run with [golang-migrate](https://github.com/golang-migrate/migrate)

### Two databases, one set of business rules

`internal/user/module.go` and `internal/ticket/module.go` both take a `port.Repository` by injection instead of building one themselves. The composition root, `cmd/api/main.go`, decides which driven adapter to wire in, based on `DB_DRIVER`:

| `DB_DRIVER` | Intended use | Adapter package | Schema |
| --- | --- | --- | --- |
| `postgres`, the default | Production | `adapter/out/postgres` | Tracked migrations in `migrations/postgres/`, applied manually with `golang-migrate` |
| `sqlite` | Local development | `adapter/out/sqlite` | Tracked migrations in `migrations/sqlite/`, embedded into the binary via `go:embed` and applied automatically on every startup, zero setup, no server required |

The SQLite driver is [`modernc.org/sqlite`](https://modernc.org/sqlite), a pure-Go implementation with no CGO requirement, so `go run ./cmd/api` works out of the box on any platform without a C toolchain. Both adapters implement the exact same `port.Repository` interface per hexagon, so the domain, use cases, and HTTP handlers never change based on which database is active. The SQL itself differs between `migrations/postgres/` and `migrations/sqlite/` because the two engines use different column types, for example `UUID` and `TIMESTAMPTZ` in PostgreSQL against `TEXT` in SQLite, but each numbered pair of files creates the same two tables with the same constraints.

The ticket hexagon additionally supports `TICKET_REPOSITORY=memory`, which forces `adapter/out/memory`, an in-memory, pre-seeded, thread-safe repository, regardless of `DB_DRIVER`. This is useful for demos or tests that should not touch disk at all.

### Teams and assignment

Every user optionally belongs to one team, tracked as a nullable `team_id` on the user. Teams themselves are a small, independent hexagon under `internal/team`, managed by administrators through the `/teams` endpoints.

When a ticket is assigned, `internal/ticket/application/usecase/ticket_service.go` checks that the user performing the assignment and the user being assigned share the same team, using an outbound port, `ticket/application/port.UserDirectory`, implemented once in `internal/ticket/adapter/out/directory` for both PostgreSQL and SQLite via `sqlx.DB.Rebind`. A user with no team can neither assign tickets nor be assigned one, and a mismatched team returns `403`. This keeps the ticket hexagon decoupled from the user hexagon's internal types: it only depends on a narrow interface it declares itself.

There is no dedicated admin bootstrap step: the very first account ever registered becomes an admin automatically, and every account after that defaults to requester. Admins manage teams and assign other users to them from the Teams screen in the frontend.

## Frontend stack

See [`ui/README.md`](ui/README.md) for the frontend's hexagonal layout. Stack: React plus TypeScript and Vite, React Router, Zustand, Axios, React Hook Form plus Zod, Ant Design, CSS Modules.

## Prerequisites

- Go 1.27 or newer
- Node.js 24 or newer, and npm
- Docker and Docker Compose, needed only for the PostgreSQL production path
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI, needed only to run PostgreSQL migrations manually or to use the optional `make migrate-sqlite-*` commands; the application itself applies SQLite migrations automatically without this CLI

## Getting started

Running CoreLog locally means running two processes at the same time: the Go API and the React frontend. Use two separate terminal windows or tabs, one for the backend and one for the frontend, and keep both running while you work.

### 1. Start the backend

Choose one of the two database options below, then start the API. Run these commands from the repository root.

#### Option A: SQLite, zero setup, recommended for local development

```sh
cp .env.example .env
```

`.env.example` already sets `DB_DRIVER=sqlite`, so no further edits to `.env` are required. Start the API:

```sh
make run
```

The API creates `./data/<DB_NAME>.db`, `./data/corelog.db` with the default `.env.example` value, and its schema automatically on first run. No Docker, no Postgres, and no migration step are required. Leave this process running: the API is now listening on `http://localhost:8080`.

#### Option B: PostgreSQL via Docker, production-like

```sh
cp .env.example .env
```

Edit `.env` and set `DB_DRIVER=postgres`, then start the database and API containers and apply the tracked migrations:

```sh
make docker-up
make migrate-up
```

The API is now listening on `http://localhost:8080`.

### 2. Start the frontend

Open a second terminal, keep the backend from step 1 running, and from the repository root run:

```sh
cd ui
cp .env.example .env
npm install
npm run dev
```

`ui/.env.example` sets `VITE_API_BASE_URL=http://localhost:8080/api/v1`, matching the backend from step 1, so no further edits are required for a default local setup. Leave this process running: the frontend dev server is now listening on `http://localhost:5173`.

### 3. Verify both are running

Open `http://localhost:5173` in a browser, register a new account, sign in, and create a ticket. The page talks to the API at `http://localhost:8080` in the background.

To check the backend on its own, without the browser:

```sh
curl http://localhost:8080/health
```

This returns `{"status":"ok"}` when the API is reachable.

## Setup scripts

`scripts/setup.sh`, for bash, and `scripts/setup.ps1`, for PowerShell, wrap the steps above into single commands. Both scripts expose the same commands and behave the same way, so pick whichever matches your shell.

| Command | Description |
| --- | --- |
| `check` | Report whether Go, Node, npm, Docker, and golang-migrate are available |
| `env` | Create `.env` and `ui/.env` from their example files if missing, without overwriting existing ones |
| `dev` | Prepare a local development environment backed by SQLite: runs `env`, downloads Go modules, installs frontend dependencies. No Docker required |
| `prod` | Prepare a production-like environment backed by PostgreSQL via Docker: verifies Docker is available, runs `env`, sets `DB_DRIVER=postgres` in a freshly created `.env`, starts the `db` and `api` containers, and applies migrations when the `migrate` CLI is installed |
| `build` | Compile the backend binary to `bin/api`, `bin/api.exe` on Windows, and build the frontend to `ui/dist` |
| `package` | Run `build`, then assemble a distributable archive under `dist/` containing the backend binary, the PostgreSQL migrations, and the frontend build, and compress it to `corelog-package.tar.gz` in bash or `corelog-package.zip` in PowerShell |
| `clean` | Remove `bin/`, `ui/dist/`, `dist/`, and the packaged archive |

Examples:

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

Running either script with no command, or an unrecognized one, prints the command list above.

## Environment variables

All variables live in `.env`, copied from `.env.example`. See [`ui/.env.example`](ui/.env.example) for the frontend's variables.

| Variable | Default | Description |
| --- | --- | --- |
| `APP_ENV` | `development` | `development` enables verbose text logs; anything else logs JSON |
| `SERVER_PORT` | `8080` | HTTP port the API listens on |
| `SERVER_READ_TIMEOUT` and `SERVER_WRITE_TIMEOUT` | `10s` | HTTP server timeouts |
| `SERVER_SHUTDOWN_TIMEOUT` | `15s` | Grace period for in-flight requests on shutdown |
| `CORS_ALLOWED_ORIGIN` | `http://localhost:5173` | Comma-separated list of origins allowed to call the API from the browser |
| `DB_DRIVER` | `sqlite` in the example file, `postgres` as the code default | `postgres` or `sqlite` |
| `DB_NAME` | `corelog` | Engine-agnostic database name. For PostgreSQL, this is the database within the server. For SQLite, this is the file `data/<DB_NAME>.db`; there is no separate SQLite-only path variable |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_SSLMODE` | `localhost`, `5432`, `corelog`, `corelog`, `disable` | PostgreSQL connection, used only when `DB_DRIVER=postgres` |
| `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME` | `25`, `25`, `5m` | PostgreSQL connection pool tuning |
| `JWT_SECRET` | `change-me-in-production` | HMAC signing secret for access tokens; set a real secret outside local dev |
| `JWT_EXPIRATION` | `24h` | Access token lifetime |
| `TICKET_REPOSITORY` | `postgres` | `postgres` follows `DB_DRIVER` despite the name, or `memory` forces the in-memory ticket adapter |

## API reference

The full contract lives in [`openapi.yaml`](openapi.yaml) at the repository root. The table below is a summary.

Base path: `/api/v1`. Endpoints listed with authentication required expect `Authorization: Bearer <token>`.

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `POST` | `/auth/register` | Not required | Create a user account. The first account ever created becomes an admin automatically; every account after that is a requester by default |
| `POST` | `/auth/login` | Not required | Authenticate and receive a JWT |
| `GET` | `/users/me` | Required | Current authenticated user |
| `PATCH` | `/users/me` | Required | Update the current user's name and email |
| `POST` | `/users/me/password` | Required | Change the current user's password, given the current one |
| `GET` | `/users/search` | Required | Search users within the caller's own team, via `?q=`, for picking a ticket assignee |
| `GET` | `/users` | Admin only | List or search every user, via `?q=` and `?team_id=` |
| `PATCH` | `/users/{id}/team` | Admin only | Assign a user to a team, or clear it with `team_id: null` |
| `GET` | `/teams` | Required | List every team |
| `POST` | `/teams` | Admin only | Create a team |
| `GET` | `/teams/{id}` | Required | Get a team by id |
| `PATCH` | `/teams/{id}` | Admin only | Rename a team |
| `DELETE` | `/teams/{id}` | Admin only | Delete a team; its members are left without a team rather than being deleted |
| `POST` | `/tickets` | Required | Create a ticket |
| `GET` | `/tickets` | Required | List tickets, filterable via `?status=`, `?priority=`, `?page=`, `?page_size=` |
| `GET` | `/tickets/{id}` | Required | Get a ticket by id |
| `PATCH` | `/tickets/{id}/status` | Required | Transition a ticket's status: `open → in_progress → resolved → closed`, with `resolved → in_progress` also allowed |
| `PATCH` | `/tickets/{id}/assign` | Required | Assign a ticket to a user; the assignee must be on the same team as the caller, or this returns `403` |
| `GET` | `/health` | Not required | Liveness check, mounted at the root path without the `/api/v1` prefix |

Ticket `status`: `open`, `in_progress`, `resolved`, `closed`. Ticket `priority`: `low`, `medium`, `high`, `critical`. Admin-only endpoints require the caller's JWT to carry the `admin` role.

## Available commands

Backend, from the `Makefile` at the repo root:

| Command | Description |
| --- | --- |
| `make run` | Run the API locally via `go run ./cmd/api` |
| `make build` | Build the API binary to `bin/api` |
| `make test` | Run Go tests |
| `make lint` | Run `go vet` |
| `make tidy` | Run `go mod tidy` |
| `make docker-up`, `make docker-down`, `make docker-logs` | Manage the PostgreSQL and API Docker Compose stack |
| `make migrate-up`, `make migrate-down` | Apply or roll back PostgreSQL migrations |
| `make migrate-create name=<name>` | Scaffold a new PostgreSQL migration pair |
| `make migrate-sqlite-up`, `make migrate-sqlite-down` | Apply or roll back SQLite migrations against `data/<DB_NAME>.db` using the `migrate` CLI directly, as an alternative to the automatic startup migration; pass `DB_NAME=<name>` to target a database other than the default `corelog` |
| `make migrate-sqlite-create name=<name>` | Scaffold a new SQLite migration pair |
| `make ui-install`, `make ui-dev`, `make ui-build`, `make ui-lint` | Proxy to the frontend's npm scripts |

Frontend, from `ui/package.json`, run inside `ui/`: `npm run dev`, `npm run build`, `npm run lint`, `npm run preview`.

## Docker

`docker-compose.yml` always starts `db`, which is PostgreSQL, and `api`. The frontend is optional and only started with the `full` profile:

```sh
docker compose up -d --build
docker compose --profile full up -d --build
```

The first command starts `db` and `api`. The second additionally starts `ui`, an nginx container serving the static production build.

### Database container: user, database, and permissions

The `db` service runs the official `postgres:17-alpine` image unmodified. That image creates the role, the database, and the ownership grant between them automatically, driven entirely by three environment variables set in `docker-compose.yml`:

| Variable in the container | Sourced from `.env` | Effect |
| --- | --- | --- |
| `POSTGRES_USER` | `DB_USER`, default `corelog` | Name of the PostgreSQL role the image creates on first startup. This role is granted the superuser privilege inside the container. |
| `POSTGRES_PASSWORD` | `DB_PASSWORD`, default `corelog` | Password set for that role. |
| `POSTGRES_DB` | `DB_NAME`, default `corelog` | Name of the database the image creates, owned by `POSTGRES_USER`. |

Because the role owns the database it creates, it already has full read and write permissions on it. No manual `CREATE ROLE`, `CREATE DATABASE`, or `GRANT` statement is required for the default setup described in this project.

This automatic setup runs only once: the first time the container starts against an empty `corelog_db_data` volume. Changing `DB_USER`, `DB_PASSWORD`, or `DB_NAME` in `.env` after that volume already holds data has no effect on the existing role or database, since the image only runs its initialization scripts against an empty data directory. To apply new credentials from a clean state, remove the volume before starting the containers again:

```sh
docker compose down
docker volume rm corelog_corelog_db_data
docker compose up -d --build
```

The volume name is the Compose project name, normally the repository directory name `corelog`, followed by an underscore and the volume name declared in `docker-compose.yml`, `corelog_db_data`. Run `docker volume ls` first to confirm the exact name on your machine, since a different directory name changes the project name prefix.

### Connecting to the database manually

Open an interactive `psql` shell inside the running container:

```sh
docker compose exec db psql -U corelog -d corelog
```

Replace `corelog` in both places with your own `DB_USER` and `DB_NAME` if you changed them from the default. From outside the container, using a local `psql` installation or a GUI client such as pgAdmin or TablePlus, connect with:

```sh
psql "host=localhost port=5432 user=corelog password=corelog dbname=corelog sslmode=disable"
```

Once connected, these `psql` meta-commands are useful for inspecting the setup:

```sql
\du
\l
\dt
```

`\du` lists roles and their attributes, `\l` lists databases, and `\dt` lists the tables in the current database.

### Creating an additional role manually

The default setup gives the application a single role with full ownership of its database, which is sufficient for local development and this project's scope. To create a second role, for example a read-only reporting account, connect as described above and run:

```sql
CREATE ROLE reporting WITH LOGIN PASSWORD 'change-me';
GRANT CONNECT ON DATABASE corelog TO reporting;
GRANT USAGE ON SCHEMA public TO reporting;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO reporting;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO reporting;
```

The last statement ensures tables created by future migrations are also readable by `reporting`, not only the tables that already existed when the grant ran.

SQLite has no equivalent user or permission model: `migrations/sqlite/` produces a single file at `data/<DB_NAME>.db`, and access to it is controlled entirely by the operating system's file permissions on that path.

## Project conventions

- No comments in source code. Names, small functions, and the hexagonal layering carry the intent; this file and `ui/README.md` are where narrative documentation belongs.
- Ports before adapters. Every use case depends on an interface declared in `application/port`, never directly on `adapter/out/*`. The composition root is the only place concrete adapters are chosen: `cmd/api/main.go` on the backend, `adapters/in/state/*.store.ts` on the frontend.
- One adapter, one persistence concern. `adapter/out/postgres` and `adapter/out/sqlite` each own their row-to-domain mapping; they never share a struct, so either can change its storage representation independently.
- Frontend imports use the `@/` alias, mapped to `ui/src`, for anything crossing a top-level boundary such as `app`, `core`, `adapters`, or `shared`. Same-directory and one-level-up imports stay relative. See `ui/README.md`.

## Troubleshooting

- `database is locked` on SQLite: only one process should hold `./data/<DB_NAME>.db` at a time. Stop any other running instance of the API.
- `migrate: command not found`: install [golang-migrate](https://github.com/golang-migrate/migrate#installation). It is required only for the PostgreSQL path, not for SQLite.
- CORS errors in the browser: make sure `CORS_ALLOWED_ORIGIN` matches the exact origin the frontend runs on, whose default is `http://localhost:5173`.
- `401 Unauthorized` on `/tickets`: the ticket and `/users/me` endpoints require `Authorization: Bearer <token>` from `POST /auth/login`.

- Módulo de autenticación (registro e inicio de sesión) probado y funcional.