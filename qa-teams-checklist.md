# QA Checklist – Módulo de Equipos (Teams)

Responsable de calidad: Jonathan Fabrizio Navarro López
Última actualización: 24/10/2026

## APF1 – Pruebas ejecutadas

### Autenticación y autorización
- [x] GET /teams sin token → esperado 401 → obtenido 401 ✅
- [x] GET /teams con token válido → esperado 200 → obtenido 200, lista vacía ✅

### Crear equipo
- [x] POST /teams con nombre válido ("QA Team") → esperado 201 → obtenido 201 ✅

### Renombrar equipo
- [x] PATCH /teams/{id} con id real → esperado 200 → obtenido 200 ✅

### Eliminar equipo
- [x] DELETE /teams/{id} con id real → esperado 204 → obtenido 204 ✅

## APF2 – Pruebas adicionales

Entorno: proyecto levantado desde cero con `docker-compose down -v` y `docker-compose up --build`. API en http://localhost:8080/api/v1. La primera cuenta registrada es admin.

### Validaciones y errores
- [x] POST /teams con nombre ya existente → esperado 409 → obtenido 409 ✅
- [x] POST /teams con nombre de 1 carácter ("A") → esperado 422 → obtenido 422 ✅
- [x] PATCH /teams/{id} renombrando a un nombre ya usado → esperado 409 → obtenido 409 ✅
- [x] GET /teams/{id} con un id existente → esperado 200 → obtenido 200 ✅
- [x] GET /teams/{id} con un id inexistente → esperado 404 → obtenido 404 ✅
- [x] DELETE /teams/{id} con un id inexistente → esperado 404 → obtenido 404 ✅
- [x] PATCH /users/{id}/team con un team_id que no es UUID → esperado 422 → obtenido 422 ✅

### Asignación de usuarios
- [x] PATCH /users/{id}/team con team_id válido (token admin) → esperado 200 → obtenido 200 ✅
- [x] DELETE /teams/{id} de un equipo con un miembro → esperado 204 → obtenido 204 ✅
- [x] El usuario del equipo eliminado sigue existiendo y queda sin equipo (team_id null) ✅

### Autorización
- [x] PATCH /users/{id}/team con token de usuario no admin → esperado 403 → obtenido 403 ✅
- [x] POST /teams con token de usuario no admin → esperado rechazo (403) → obtenido ✅

### Regresión tras los cambios del APF2 (PR #12 y #13)
- [x] docker-compose define los servicios db, migrate, api y ui ✅
- [x] Tras docker-compose up: db, api y ui en Up; migrate en Exited (0) ✅
- [x] Registro de usuario desde cero, sin migraciones manuales → esperado 201 → obtenido 201 ✅
- [ ] El frontend carga en http://localhost:5173 → obtenido 201
- [x] POST /teams/{uuid} → esperado 404 o 405 → obtenido ✅
- [ ] POST /teams/{id} con el {id} literal (caso original de la incidencia) → esperado 404 o 405 → obtenido 

### Búsqueda de tickets por título (PR #12)
- [x] Con 3 tickets ("Error de login", "Falla en reportes", "Login lento"), buscar "login" → esperado 2 resultados → obtenido ✅
- [x] Buscar "LOGIN" en mayúsculas → esperado: ver si distingue mayúsculas → obtenido ✅
- [x] Buscar "zzz" → esperado lista vacía o mensaje de sin resultados → obtenido ✅
- [x] Borrar el texto de búsqueda → esperado: vuelven los 3 tickets → obtenido ✅

## Revisiones de Pull Requests realizadas

| PR | Rama | Qué se probó | Resultado |
|----|------|--------------|-----------|
| #13 | fix/docker-compose-completo | Arranque desde cero: servicios db, migrate, api y ui; migraciones automáticas; registro | ✅ Funciona |
| #12 | feature/tickets-busqueda | Búsqueda de tickets por título en la interfaz | Pendiente |

## Incidencias

1. **Dockerfile.api no copiaba la carpeta `migrations` en la etapa de build**, lo que rompía la compilación del backend. Corregido en el PR #9.
2. **docker-compose.yml no incluía el servicio del frontend.** Corregido en el PR #13 y verificado.
3. **Las migraciones no se ejecutaban automáticamente** (error 500 en el primer registro). Corregido en el PR #13 con el servicio `migrate` y verificado.
4. **Posible bug de enrutamiento:** un POST a `/teams/{id}` creó un equipo en vez de devolver 404/405. Pendiente de confirmar con el caso original exacto.