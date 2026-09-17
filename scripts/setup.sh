#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
UI_DIR="$ROOT_DIR/ui"
BIN_DIR="$ROOT_DIR/bin"
DIST_DIR="$ROOT_DIR/dist"

log() {
  printf '[%s] %s\n' "$1" "$2"
}

fail() {
  printf '[error] %s\n' "$1" >&2
  exit 1
}

has_command() {
  command -v "$1" >/dev/null 2>&1
}

check_go() {
  if has_command go; then
    log check "Go found: $(go version)"
    return 0
  fi
  log check "Go not found"
  return 1
}

check_node() {
  if has_command node && has_command npm; then
    log check "Node found: $(node --version), npm found: $(npm --version)"
    return 0
  fi
  log check "Node or npm not found"
  return 1
}

check_docker() {
  if ! has_command docker; then
    log check "Docker not found"
    return 1
  fi
  if ! docker info >/dev/null 2>&1; then
    log check "Docker found but the daemon is not reachable: $(docker --version)"
    return 1
  fi
  if docker compose version >/dev/null 2>&1; then
    log check "Docker found: $(docker --version), Docker Compose found: $(docker compose version --short 2>/dev/null || echo unknown)"
    return 0
  fi
  log check "Docker found: $(docker --version), Docker Compose plugin not found"
  return 1
}

check_migrate() {
  if has_command migrate; then
    log check "golang-migrate found: $(migrate -version 2>&1)"
    return 0
  fi
  log check "golang-migrate not found, only required for manual PostgreSQL or SQLite migration commands"
  return 1
}

cmd_check() {
  log check "Checking prerequisites for CoreLog"
  check_go || true
  check_node || true
  check_docker || true
  check_migrate || true
}

create_env_file() {
  local example_path="$1"
  local target_path="$2"
  local label="$3"
  if [ -f "$target_path" ]; then
    log env "$label already exists, leaving it untouched: $target_path"
    return 0
  fi
  if [ ! -f "$example_path" ]; then
    fail "$label example file is missing: $example_path"
  fi
  cp "$example_path" "$target_path"
  log env "$label created from example: $target_path"
}

cmd_env() {
  create_env_file "$ROOT_DIR/.env.example" "$ROOT_DIR/.env" "Backend .env"
  create_env_file "$UI_DIR/.env.example" "$UI_DIR/.env" "Frontend .env"
}

cmd_dev() {
  log dev "Preparing the local development environment, SQLite backend, no Docker required"
  check_go || fail "Go is required for local development, install it and re-run this command"
  cmd_env
  log dev "Downloading Go module dependencies"
  (cd "$ROOT_DIR" && go mod download)
  if check_node; then
    log dev "Installing frontend dependencies"
    (cd "$UI_DIR" && npm install)
  else
    log dev "Skipping frontend dependency install, Node or npm not found"
  fi
  log dev "Environment ready. Run 'make run' in one terminal and 'cd ui && npm run dev' in another"
}

cmd_prod() {
  log prod "Preparing the production-like environment, PostgreSQL via Docker"
  check_docker || fail "Docker and the Docker Compose plugin are required for the production path, install them and re-run this command"
  local created_env=0
  if [ ! -f "$ROOT_DIR/.env" ]; then
    created_env=1
  fi
  cmd_env
  if [ "$created_env" -eq 1 ]; then
    sed -i.bak 's/^DB_DRIVER=.*/DB_DRIVER=postgres/' "$ROOT_DIR/.env" && rm -f "$ROOT_DIR/.env.bak"
    log prod "Set DB_DRIVER=postgres in the newly created .env"
  else
    log prod "Existing .env left untouched, confirm DB_DRIVER=postgres yourself if needed"
  fi
  log prod "Starting the db and api containers"
  (cd "$ROOT_DIR" && docker compose up -d --build)
  if has_command migrate; then
    log prod "Applying PostgreSQL migrations"
    local db_url="postgres://corelog:corelog@localhost:5432/corelog?sslmode=disable"
    (cd "$ROOT_DIR" && migrate -path migrations/postgres -database "$db_url" up) || log prod "Migration command failed, check DB_URL matches your .env and retry with 'make migrate-up'"
  else
    log prod "golang-migrate not found, run 'make migrate-up' after installing it from https://github.com/golang-migrate/migrate"
  fi
  log prod "Environment ready. The API is listening on http://localhost:8080"
}

cmd_build() {
  check_go || fail "Go is required to build the backend"
  local goexe
  goexe="$(go env GOEXE)"
  mkdir -p "$BIN_DIR"
  log build "Compiling the backend binary"
  (cd "$ROOT_DIR" && go build -o "bin/api${goexe}" ./cmd/api)
  log build "Backend binary written to bin/api${goexe}"
  if check_node; then
    log build "Building the frontend"
    (cd "$UI_DIR" && npm install && npm run build)
    log build "Frontend build written to ui/dist"
  else
    log build "Skipping frontend build, Node or npm not found"
  fi
}

cmd_package() {
  cmd_build
  local goexe
  goexe="$(go env GOEXE)"
  rm -rf "$DIST_DIR"
  mkdir -p "$DIST_DIR/migrations"
  cp "$BIN_DIR/api${goexe}" "$DIST_DIR/"
  cp -r "$ROOT_DIR/migrations/postgres" "$DIST_DIR/migrations/postgres"
  cp "$ROOT_DIR/.env.example" "$DIST_DIR/.env.example"
  if [ -d "$UI_DIR/dist" ]; then
    cp -r "$UI_DIR/dist" "$DIST_DIR/ui"
  fi
  local archive_path="$ROOT_DIR/corelog-package.tar.gz"
  (cd "$DIST_DIR" && tar czf "$archive_path" .)
  log package "Package assembled at dist/ and archived to $(basename "$archive_path")"
}

cmd_clean() {
  rm -rf "$BIN_DIR" "$UI_DIR/dist" "$DIST_DIR" "$ROOT_DIR/corelog-package.tar.gz"
  log clean "Removed bin, ui/dist, dist, and the packaged archive"
}

usage() {
  cat <<'EOF'
Usage: scripts/setup.sh <command>

Commands:
  check     Report whether Go, Node, npm, Docker, and golang-migrate are available
  env       Create .env and ui/.env from their example files if missing
  dev       Prepare a local development environment backed by SQLite, no Docker required
  prod      Prepare a production-like environment backed by PostgreSQL via Docker
  build     Compile the backend binary and build the frontend
  package   Build and assemble a distributable archive under dist/
  clean     Remove build and package output
EOF
}

main() {
  local command="${1:-}"
  case "$command" in
    check) cmd_check ;;
    env) cmd_env ;;
    dev) cmd_dev ;;
    prod) cmd_prod ;;
    build) cmd_build ;;
    package) cmd_package ;;
    clean) cmd_clean ;;
    *) usage ;;
  esac
}

main "$@"
