[CmdletBinding()]
param(
    [Parameter(Position = 0)]
    [string]$Command = ""
)

$ErrorActionPreference = "Continue"

$RootDir = Split-Path -Parent $PSScriptRoot
$UiDir = Join-Path $RootDir "ui"
$BinDir = Join-Path $RootDir "bin"
$DistDir = Join-Path $RootDir "dist"

function Write-Log {
    param([string]$Stage, [string]$Message)
    Write-Host "[$Stage] $Message"
}

function Write-Failure {
    param([string]$Message)
    Write-Host "[error] $Message" -ForegroundColor Red
    exit 1
}

function Assert-LastExitCode {
    param([string]$Message)
    if ($LASTEXITCODE -ne 0) {
        Write-Failure $Message
    }
}

function Test-CommandAvailable {
    param([string]$Name)
    return [bool](Get-Command $Name -ErrorAction SilentlyContinue)
}

function Test-Go {
    if (Test-CommandAvailable "go") {
        Write-Log "check" "Go found: $(go version)"
        return $true
    }
    Write-Log "check" "Go not found"
    return $false
}

function Test-Node {
    if ((Test-CommandAvailable "node") -and (Test-CommandAvailable "npm")) {
        Write-Log "check" "Node found: $(node --version), npm found: $(npm --version)"
        return $true
    }
    Write-Log "check" "Node or npm not found"
    return $false
}

function Test-Docker {
    if (-not (Test-CommandAvailable "docker")) {
        Write-Log "check" "Docker not found"
        return $false
    }
    docker info *> $null
    if ($LASTEXITCODE -ne 0) {
        Write-Log "check" "Docker found but the daemon is not reachable: $(docker --version)"
        return $false
    }
    docker compose version *> $null
    if ($LASTEXITCODE -eq 0) {
        Write-Log "check" "Docker found: $(docker --version), Docker Compose found"
        return $true
    }
    Write-Log "check" "Docker found: $(docker --version), Docker Compose plugin not found"
    return $false
}

function Test-Migrate {
    if (Test-CommandAvailable "migrate") {
        Write-Log "check" "golang-migrate found: $(migrate -version 2>&1)"
        return $true
    }
    Write-Log "check" "golang-migrate not found, only required for manual PostgreSQL or SQLite migration commands"
    return $false
}

function Invoke-Check {
    Write-Log "check" "Checking prerequisites for CoreLog"
    Test-Go | Out-Null
    Test-Node | Out-Null
    Test-Docker | Out-Null
    Test-Migrate | Out-Null
}

function New-EnvFile {
    param([string]$ExamplePath, [string]$TargetPath, [string]$Label)
    if (Test-Path $TargetPath) {
        Write-Log "env" "$Label already exists, leaving it untouched: $TargetPath"
        return
    }
    if (-not (Test-Path $ExamplePath)) {
        Write-Failure "$Label example file is missing: $ExamplePath"
    }
    Copy-Item $ExamplePath $TargetPath
    Write-Log "env" "$Label created from example: $TargetPath"
}

function Invoke-Env {
    New-EnvFile (Join-Path $RootDir ".env.example") (Join-Path $RootDir ".env") "Backend .env"
    New-EnvFile (Join-Path $UiDir ".env.example") (Join-Path $UiDir ".env") "Frontend .env"
}

function Invoke-Dev {
    Write-Log "dev" "Preparing the local development environment, SQLite backend, no Docker required"
    if (-not (Test-Go)) {
        Write-Failure "Go is required for local development, install it and re-run this command"
    }
    Invoke-Env
    Write-Log "dev" "Downloading Go module dependencies"
    Push-Location $RootDir
    go mod download
    Assert-LastExitCode "go mod download failed"
    Pop-Location
    if (Test-Node) {
        Write-Log "dev" "Installing frontend dependencies"
        Push-Location $UiDir
        npm install
        Assert-LastExitCode "npm install failed"
        Pop-Location
    }
    else {
        Write-Log "dev" "Skipping frontend dependency install, Node or npm not found"
    }
    Write-Log "dev" "Environment ready. Run 'go run ./cmd/api' in one terminal and 'cd ui; npm run dev' in another"
}

function Invoke-Prod {
    Write-Log "prod" "Preparing the production-like environment, PostgreSQL via Docker"
    if (-not (Test-Docker)) {
        Write-Failure "Docker and the Docker Compose plugin are required for the production path, install them and re-run this command"
    }
    $envPath = Join-Path $RootDir ".env"
    $createdEnv = -not (Test-Path $envPath)
    Invoke-Env
    if ($createdEnv) {
        (Get-Content $envPath) -replace '^DB_DRIVER=.*', 'DB_DRIVER=postgres' | Set-Content $envPath
        Write-Log "prod" "Set DB_DRIVER=postgres in the newly created .env"
    }
    else {
        Write-Log "prod" "Existing .env left untouched, confirm DB_DRIVER=postgres yourself if needed"
    }
    Write-Log "prod" "Starting the db and api containers"
    Push-Location $RootDir
    docker compose up -d --build
    Assert-LastExitCode "docker compose up failed"
    Pop-Location
    if (Test-CommandAvailable "migrate") {
        Write-Log "prod" "Applying PostgreSQL migrations"
        $dbUrl = "postgres://corelog:corelog@localhost:5432/corelog?sslmode=disable"
        Push-Location $RootDir
        migrate -path migrations/postgres -database $dbUrl up
        if ($LASTEXITCODE -ne 0) {
            Write-Log "prod" "Migration command failed, check DB_URL matches your .env and retry with 'make migrate-up'"
        }
        Pop-Location
    }
    else {
        Write-Log "prod" "golang-migrate not found, run 'make migrate-up' after installing it from https://github.com/golang-migrate/migrate"
    }
    Write-Log "prod" "Environment ready. The API is listening on http://localhost:8080"
}

function Invoke-Build {
    if (-not (Test-Go)) {
        Write-Failure "Go is required to build the backend"
    }
    $goExe = (go env GOEXE).Trim()
    New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
    Write-Log "build" "Compiling the backend binary"
    Push-Location $RootDir
    go build -o "bin/api$goExe" ./cmd/api
    Assert-LastExitCode "go build failed"
    Pop-Location
    Write-Log "build" "Backend binary written to bin/api$goExe"
    if (Test-Node) {
        Write-Log "build" "Building the frontend"
        Push-Location $UiDir
        npm install
        Assert-LastExitCode "npm install failed"
        npm run build
        Assert-LastExitCode "npm run build failed"
        Pop-Location
        Write-Log "build" "Frontend build written to ui/dist"
    }
    else {
        Write-Log "build" "Skipping frontend build, Node or npm not found"
    }
}

function Invoke-Package {
    Invoke-Build
    $goExe = (go env GOEXE).Trim()
    if (Test-Path $DistDir) {
        Remove-Item -Recurse -Force $DistDir
    }
    New-Item -ItemType Directory -Force -Path (Join-Path $DistDir "migrations") | Out-Null
    Copy-Item (Join-Path $BinDir "api$goExe") $DistDir
    Copy-Item -Recurse (Join-Path $RootDir "migrations/postgres") (Join-Path $DistDir "migrations/postgres")
    Copy-Item (Join-Path $RootDir ".env.example") $DistDir
    $uiDist = Join-Path $UiDir "dist"
    if (Test-Path $uiDist) {
        Copy-Item -Recurse $uiDist (Join-Path $DistDir "ui")
    }
    $archivePath = Join-Path $RootDir "corelog-package.zip"
    if (Test-Path $archivePath) {
        Remove-Item -Force $archivePath
    }
    Compress-Archive -Path (Join-Path $DistDir "*") -DestinationPath $archivePath
    Write-Log "package" "Package assembled at dist/ and archived to $(Split-Path -Leaf $archivePath)"
}

function Invoke-Clean {
    Remove-Item -Recurse -Force -ErrorAction SilentlyContinue $BinDir
    Remove-Item -Recurse -Force -ErrorAction SilentlyContinue (Join-Path $UiDir "dist")
    Remove-Item -Recurse -Force -ErrorAction SilentlyContinue $DistDir
    Remove-Item -Force -ErrorAction SilentlyContinue (Join-Path $RootDir "corelog-package.zip")
    Write-Log "clean" "Removed bin, ui/dist, dist, and the packaged archive"
}

function Show-Usage {
    Write-Host "Usage: scripts/setup.ps1 <command>"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  check     Report whether Go, Node, npm, Docker, and golang-migrate are available"
    Write-Host "  env       Create .env and ui/.env from their example files if missing"
    Write-Host "  dev       Prepare a local development environment backed by SQLite, no Docker required"
    Write-Host "  prod      Prepare a production-like environment backed by PostgreSQL via Docker"
    Write-Host "  build     Compile the backend binary and build the frontend"
    Write-Host "  package   Build and assemble a distributable archive under dist/"
    Write-Host "  clean     Remove build and package output"
}

switch ($Command) {
    "check" { Invoke-Check }
    "env" { Invoke-Env }
    "dev" { Invoke-Dev }
    "prod" { Invoke-Prod }
    "build" { Invoke-Build }
    "package" { Invoke-Package }
    "clean" { Invoke-Clean }
    default { Show-Usage }
}
