![Go Version](https://img.shields.io/badge/go-1.24+-darkgreen)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-12+-blue)
![License](https://img.shields.io/github/license/your-repo/weather-api)


# Weather API

A RESTful API service for retrieving, storing, and managing weather data from OpenWeatherMap.

## Table of Contents

- [Features](#features)
- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Quickstart](#quickstart)
- [Environment Variables](#environment-variables)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [Docker Setup](#docker-setup)
- [API Documentation](#api-documentation)
  - [Health Check Endpoints](#health-check-endpoints)
  - [Swagger UI](#swagger-ui)
  - [API Endpoints](#api-endpoints)
  - [Example Requests](#example-requests)
  - [Postman Collection](#postman-collection)
- [Authentication (JWT)](#authentication-jwt)
- [Caching Strategy](#caching-strategy)
- [Error Handling](#error-handling)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Makefile Usage](#makefile-usage)
- [Troubleshooting](#troubleshooting)
- [License](#license)

## Features

- Fetch and store current weather data by city and country
- Cache weather data in Redis to reduce external API calls
- Complete CRUD operations for weather records
- Input validation with detailed error messages
- Swagger API documentation
- Comprehensive error handling

## Technology Stack

- **Go**: Core programming language (1.24+)
- **Gin**: Web framework for routing and middleware
- **PostgreSQL**: Primary database for storing weather data
- **Redis**: Caching layer for weather data
- **Swagger**: API documentation
- **go-playground/validator**: Request validation

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher
- OpenWeatherMap API key (obtain at [OpenWeather](https://openweathermap.org/api))

## Setup

## Quickstart
```bash
# Clone and enter the project
git clone https://github.com/OmidRasouli/weather-api weather-api
cd weather-api

# Prepare environment
cp .env.example .env
# Edit .env and set OPENWEATHER_API_KEY and DB/Redis settings as needed

# Run everything via Docker
docker-compose up -d

# Verify health
curl -s http://localhost:8080/health
```

### Environment Variables

Create a `.env` file in the project root with the following variables:

```properties
# Server Configuration
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=weather
DB_SSLMODE=disable

# OpenWeatherMap API
OPENWEATHER_API_KEY=your_api_key_here

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_TTL=600
```

Configuration reference:
- SERVER_PORT: API server port (default 8080)
- DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE: PostgreSQL connection params
- OPENWEATHER_API_KEY: Your OpenWeather API key (required)
- REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, REDIS_DB: Redis connection params
- REDIS_TTL: Cache TTL in seconds (default 600)

### Database Setup

1. Create a PostgreSQL database named `weather`
2. The application will automatically run migrations on startup

### Running the Application

#### Option 1: Local Development
```bash
# Get dependencies
go mod download

# Run the application
go run cmd/main.go
```

#### Option 2: Docker (Recommended)
```bash
# Copy environment file and configure
cp .env
# Edit .env and add your OPENWEATHER_API_KEY

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f weather-api

# Stop services
docker-compose down
```

The server will start on the configured port (default: 8080).

## Docker Setup

### Prerequisites for Docker
- Docker
- Docker Compose
- OpenWeatherMap API key

### Docker Services
The application runs with three services:
- **weather-api**: The Go application
- **postgres**: PostgreSQL database
- **redis**: Redis cache

### Docker Commands
```bash
# Build and start all services
docker-compose up -d

# View logs for specific service
docker-compose logs weather-api
docker-compose logs postgres
docker-compose logs redis

# Rebuild the application
docker-compose build weather-api

# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: This deletes all data)
docker-compose down -v
```

### Environment Configuration for Docker
Create a `.env` file in the project root:
```bash
cp .env.example .env
```
Then edit the `.env` file and set your `OPENWEATHER_API_KEY`.

## API Documentation

### Health Check Endpoints

The API provides health check endpoints to monitor service status:

| Endpoint | Description |
|----------|-------------|
| GET /health | Basic health check that returns 200 OK if the service is running |
| GET /health/ready | Readiness check that verifies connections to PostgreSQL and Redis |
| GET /health/live | Liveness check for container orchestration systems like Kubernetes |

Response format:
```json
{
  "status": "UP",
  "components": {
    "database": "UP",
    "redis": "UP",
    "api": "UP"
  },
  "version": "1.0.0"
}
```

Status codes:
- 200: Service is healthy
- 503: Service is unhealthy or dependencies are unavailable

### Swagger UI

Once the application is running, access Swagger documentation at:
```
http://localhost:8080/swagger/index.html
```

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /weather | List all weather records |
| GET | /weather/:id | Get weather by ID |
| POST | /weather | Fetch and store weather for a city/country |
| PUT | /weather/:id | Update a weather record |
| DELETE | /weather/:id | Delete a weather record |
| GET | /weather/latest/:city | Get latest weather for a city |

### Example Requests

#### Fetch Weather Data
```bash
curl -X POST http://localhost:8080/weather \
  -H "Content-Type: application/json" \
  -d '{"city": "London", "country": "GB"}'
```

#### Get Latest Weather for a City
```bash
curl -X GET http://localhost:8080/weather/latest/London
```

### Postman Collection

This project includes a Postman collection for easier API testing:

1. **Import the Collection and Environment**:
   - Files are located in the `/postman` directory
   - Import `Weather_API.postman_collection.json` into Postman
   - Import `Weather_API_Environment.postman_environment.json`

2. **Configure Environment**:
   - Select the "Weather API Environment" from the environment dropdown
   - Verify that `baseUrl` is set correctly (default: `http://localhost:8080`)
   - Other variables like `weatherId` will be populated automatically by test scripts

3. **Using the Collection**:
   - The collection includes requests for all API endpoints
   - Start with "Fetch and Store Weather" to create a record
   - Test scripts will automatically extract and store the created weather ID
   - Use other requests to test CRUD operations with the saved ID

4. **Custom Variables**:
   - To test with different cities, update the `cityName` variable
   - For manual testing, update the `weatherId` variable with a valid UUID

## Authentication (JWT)

- Set env: `JWT_SECRET`, `ADMIN_USERNAME`, `ADMIN_PASSWORD`.
- Obtain a token:
  ```
  POST /login
  {
    "username": "admin",
    "password": "strong-password"
  }
  ```
- Use the token:
  ```
  Authorization: Bearer <token>
  ```

Protected endpoints:
- `POST /weather`
- `PUT /weather/:id`
- `DELETE /weather/:id`

Public endpoints remain:
- `GET /weather`
- `GET /weather/:id`
- `GET /weather/latest/:city`

## Caching Strategy

Weather data is cached in Redis with the following approach:

- Weather data is cached by city and country key (e.g., `weather:London:UK`)
- Default TTL is 10 minutes (configurable via `REDIS_TTL`)
- Cache is invalidated automatically when TTL expires
- Cache hits reduce load on the OpenWeatherMap API

## Error Handling

The API provides consistent error responses with appropriate HTTP status codes:

- 400: Bad Request (validation errors)
- 404: Not Found
- 500: Internal Server Error
- 502: Bad Gateway (external API errors)

Error responses include detailed messages to help diagnose issues.

## Project Structure

The project follows clean architecture principles:

```
weather-api/
├── cmd/                   # Application entry points
├── docs/                  # Swagger documentation
├── infrastructure/        # External systems interfaces (DB, Redis)
├── internal/
│   ├── application/       # Application services
│   ├── configs/           # Configuration management
│   ├── domain/            # Domain models and interfaces
│   ├── infrastructure/    # Infrastructure implementations
│   └── interfaces/        # API controllers and routes
├── pkg/                   # Reusable packages
└── scripts/               # Utility scripts
```

## Testing

Run tests with:

```bash
go test ./...
```

Or use the Makefile targets:

```bash
make           # default: runs all tests (integration then unit)
make test-unit
make test-integration
make test-coverage     # prints summary and writes coverage.out
make coverage-html     # generates coverage.html from coverage.out
make test-race         # runs tests with -race
make test-fresh        # clears test cache, then runs tests
make build
make run
make deps
make clean
make health            # quick /health check (requires server running)
make test-watch        # watches integration tests (requires 'entr')
```

## Makefile Usage

Common workflows:
- Quick test run: make or make test
- Unit tests only: make test-unit
- Integration tests only: make test-integration
- Coverage report: make test-coverage then make coverage-html
- Data race checks: make test-race
- Fresh (no cache): make test-fresh
- Build/run: make build, make run

## Troubleshooting
- 401/403 from OpenWeather: ensure OPENWEATHER_API_KEY is set and valid.
- Cannot connect to DB/Redis: verify DB_* and REDIS_* envs and that containers are up (docker-compose ps).
- Swagger not reachable: confirm server is running on SERVER_PORT and visit /swagger/index.html.
- Go version mismatch: align the README badge with your go.mod Go version.
- Tests failing due to cache: use make test-fresh.

## License

[MIT](LICENSE)