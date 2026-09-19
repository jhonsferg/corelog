# QA Checklist – Módulo de Equipos (Teams)

Responsable: Jonathan Fabrizio Navarro López
Fecha: 18/09/2026

## Autenticación y autorización
- [x] GET /teams sin token → esperado 401 → obtenido 401 ✅
- [x] GET /teams con token válido → esperado 200 → obtenido 200, lista vacía ✅

## Crear equipo
- [x] POST /teams con nombre válido ("QA Team") → esperado 201 → obtenido 201 ✅

## Renombrar equipo
- [x] PATCH /teams/{id} con id real → esperado 200 → obtenido 200 ✅

## Eliminar equipo
- [x] DELETE /teams/{id} con id real → esperado 204 → obtenido 204 ✅

## Incidencias encontradas durante las pruebas

1. **Dockerfile.api no copiaba la carpeta `migrations` en la etapa de build**, lo que rompía la compilación del backend en Docker. Corregido agregando `COPY migrations ./migrations` antes del `RUN go build`.

2. **`docker-compose.yml` no incluye un servicio para el frontend (React/UI)**, solo `api` y `db`. Actualmente hay que correr el frontend por separado con `npm run dev`.

3. **Las migraciones de la base de datos no se ejecutan automáticamente** al levantar los contenedores con `docker-compose up`, causando error 500 ("relation users does not exist") en el primer registro de usuario. Hay que correrlas manualmente con `make migrate-up` (o con un contenedor `migrate/migrate` si no se tiene el CLI instalado).

4. **Posible bug de enrutamiento:** un `POST` a `/teams/{id}` (método no definido en el openapi.yaml para esa ruta) no devuelve `404`/`405` como se esperaría, sino que crea un nuevo equipo, ignorando el `{id}` del path. Pendiente de confirmar con el responsable del backend.