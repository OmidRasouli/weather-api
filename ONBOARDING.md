# Weather API – Onboarding

## 1) Overview
- Purpose: RESTful API for weather data management. Fetches weather data from OpenWeather, persists in Postgres, and caches in Redis. Exposes HTTP endpoints with Swagger docs.
- Tech stack:
  - Go, Gin (HTTP server/router)
  - Postgres (primary storage)
  - Redis (optional cache)
  - Swagger (swaggo) for API docs
  - Internal migration manager
  - Structured logging and validation utilities

## 2) Folder structure
- cmd/
  - main.go – Application entrypoint; wiring, startup, shutdown.
- config/ – Configuration loading (env/files).
- docs/ – Generated Swagger docs (served at /swagger/*any).
- infrastructure/
  - database/
    - cache/ – Redis client wiring.
    - database/ – Postgres connection wiring.
- internal/
  - application/
    - auth/ – Auth use cases.
    - interfaces/ – Cross-layer interfaces (Database, Cache, etc.).
    - service/ – Application services (e.g., WeatherService).
  - database/
    - migrations/ – Migration instance/manager and migration files.
  - domain/
    - services/ – Domain services (e.g., AuthService).
  - infrastructure/
    - database/
      - postgres/
        - weather/ – Weather repository implementation (Postgres).
    - openweather/ – OpenWeather API client.
  - interfaces/
    - http/
      - controller/ – HTTP controllers/handlers.
      - routers/ – Router setup and middleware wiring.
- pkg/
  - logger/ – Logger initialization and helpers.
  - validator/ – Validator initialization and custom rules.

## 3) Key modules and responsibilities
- main.go: Initializes logger and validator; loads config; sets up DB/Redis; runs migrations; wires repositories, clients, services, controllers, and router; starts HTTP server.
- config: Exposes Load() to read configuration and environment variables.
- database/database (postgres): Creates and configures DB connection (pool sizing, lifetime, SSL).
- database/cache (redis): Creates Redis client. Non-fatal if unavailable.
- internal/database/migrations: Creates migration instance and runs pending migrations on startup.
- internal/infrastructure/openweather: Thin client for OpenWeather API.
- internal/infrastructure/database/postgres/weather: Weather data repository (Postgres-specific).
- internal/application/service: WeatherService orchestrating repository + external API + caching.
- internal/application/auth + internal/domain/services: Auth use cases and domain logic.
- internal/interfaces/http/controller: HTTP layer adapters for services.
- internal/interfaces/http/routers: Gin router setup, routes, and middleware.
- pkg/logger: Centralized structured logging.
- pkg/validator: Validation rules and initialization.

## 4) Main execution flow
1. Logger initialized; config loaded.
2. Validator initialized early.
3. RunDatabase:
   - Connects to Postgres.
   - Connects to Redis (optional; logs warning if it fails).
   - Builds migration instance and runs migrations (fail-fast on errors).
4. RunServer:
   - Construct Weather repository (Postgres) and OpenWeather client.
   - Wire WeatherService with repo + API + Redis.
   - Build controllers and router; mount Swagger at /swagger/*any.
   - Start Gin server on :<Server.Port>.
5. Deferred cleanup closes DB and Redis on shutdown.

## 5) Configuration and environment variables
- Handled by config.Load(). It reads structured config for:
  - Server: Port.
  - OpenWeather: APIKey.
  - Database: Host, Port, User, Password, DBName, SSLMode, MaxIdleConns, MaxOpenConns, ConnMaxLifetime.
  - Redis: Connection settings (see cache package).
- Notes:
  - OpenWeather API key is required for live data fetching.
  - Redis is optional; the app continues without cache if Redis is down.
  - Migrations run automatically at startup.
- Tip: See config package for the exact variable names/sources. Common environment keys resemble:
  - SERVER_PORT, OPENWEATHER_API_KEY
  - DATABASE_HOST, DATABASE_PORT, DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME, DATABASE_SSLMODE, DATABASE_MAX_IDLE_CONNS, DATABASE_MAX_OPEN_CONNS, DATABASE_CONN_MAX_LIFETIME
  - REDIS_HOST/PORT/PASSWORD/DB or REDIS_ADDR

## 6) Local setup and run
- Prerequisites:
  - Go installed
  - Postgres running and reachable
  - (Optional) Redis running
  - OpenWeather API key
- Steps:
  - Set environment variables per config package.
  - Ensure the Postgres database exists and is reachable.
  - Run:
    - go run ./cmd
  - Swagger UI:
    - http://localhost:<Server.Port>/swagger/index.html

- Swagger generation (if you edit annotations):
  - Install CLI: go install github.com/swaggo/swag/cmd/swag@latest
  - From repo root: swag init -g cmd/main.go -o docs

## 7) Running tests
- Unit tests:
  - go test ./...
- Integration tests (if present) may require Postgres/Redis. Ensure services are running and test config points to test DBs.

## 8) Development and deployment notes
- Migrations: Auto-run at startup; failures are fatal. Keep migrations idempotent and committed.
- Redis: Non-blocking optional dependency; watch logs for cache connectivity warnings.
- DB pooling: Tuned via MaxIdleConns, MaxOpenConns, ConnMaxLifetime in config.
- Secrets: Provide via environment or secret manager; do not commit.
- Observability: Uses centralized logger; prefer structured logs.
- Swagger: Keep annotations current; regenerate on API changes.
- Production:
  - Set SSLMode appropriately for Postgres.
  - Use a reverse proxy/TLS terminator as needed.
  - Configure health checks and resource limits for Postgres/Redis dependencies.
