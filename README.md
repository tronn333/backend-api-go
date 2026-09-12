# Marketplace Backend API (Go)

A REST API for a simple marketplace built with Go, Gin, and PostgreSQL. It provides
authentication (JWT), user profiles, product management, and purchase/order history.

## Tech stack

- Go 1.21+
- [Gin](https://github.com/gin-gonic/gin) web framework
- PostgreSQL (via `lib/pq`)
- JWT (HS256) authentication, bcrypt password hashing

## Endpoints

| Method | Path                    | Auth   | Description                     |
| ------ | ----------------------- | ------ | ------------------------------- |
| GET    | `/health`               | public | Health check                    |
| POST   | `/auth/register`        | public | Create an account               |
| POST   | `/auth/login`           | public | Log in (returns a JWT)          |
| GET    | `/products`             | public | List products (paginated)       |
| GET    | `/products/:id`         | public | Get a single product            |
| GET    | `/users/me`             | JWT    | Current user profile            |
| PATCH  | `/users/me`             | JWT    | Update username / email         |
| DELETE | `/users/me`             | JWT    | Delete account                  |
| POST   | `/products`             | JWT    | Create a product                |
| PATCH  | `/products/:id`         | JWT    | Update own product (admin: any) |
| DELETE | `/products/:id`         | JWT    | Delete own product (admin: any) |
| POST   | `/purchases`            | JWT    | Purchase a product              |
| GET    | `/purchases`            | JWT    | List own purchases (paginated)  |
| GET    | `/purchases/:id`        | JWT    | Get a single purchase           |
| PATCH  | `/purchases/:id/cancel` | JWT    | Cancel and refund stock         |

Protected routes require the header `Authorization: Bearer <token>`.

## Environment variables

| Variable      | Default               | Description                   |
| ------------- | --------------------- | ----------------------------- |
| `DB_HOST`     | `localhost`           | PostgreSQL host               |
| `DB_PORT`     | `5432`                | PostgreSQL port               |
| `DB_NAME`     | `marketplace_db`      | Database name                 |
| `DB_USER`     | `postgres`            | Database user                 |
| `DB_PASSWORD` | `marketplace123`      | Database password             |
| `JWT_SECRET`  | `changeme` (dev only) | HMAC secret used to sign JWTs |
| `SERVER_PORT` | `8080`                | HTTP port                     |

Copy `.env.example` to `.env` and adjust as needed. Docker Compose overrides
`DB_HOST` to `db` automatically.

## Run locally (without Docker)

```bash
# 1. Create and start a PostgreSQL database, then install deps
go mod download

# 2. Configure env (or export the variables)
cp .env.example .env

# 3. Run (applies migrations automatically on startup)
go run ./cmd/api
```

Database tables and indexes are created automatically via `db.Migrate` on startup.

## Run with Docker Compose

Requires Docker and Compose v2+ (`docker compose` or the standalone `docker-compose`).

```bash
# Copy and adjust the environment file (optional — sane defaults are provided)
cp .env.example .env

# Build and start both the API and PostgreSQL
docker-compose up -d --build

# Check status / logs
docker-compose ps
docker-compose logs -f api

# Stop
docker-compose down          # keep data
docker-compose down -v       # remove the database volume too
```

- The API is exposed at `http://localhost:8080`.
- PostgreSQL is exposed on host port `5433` (`postgresql://postgres:marketplace123@localhost:5433/marketplace_db`)
  so it doesn't clash with a local Postgres on `5432`. Change `DB_EXPOSE_PORT` if needed.
- The `api` service waits for the database to be healthy before starting, and the app
  runs its migrations on first boot. On `docker-compose down`/`Ctrl+C` the process
  receives SIGTERM and shuts down gracefully, draining in-flight requests.

## Build the image on its own

```bash
docker build -t marketplace-api .
docker run --rm -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_NAME=marketplace_db \
  -e DB_USER=postgres \
  -e DB_PASSWORD=marketplace123 \
  -e JWT_SECRET=super-secret-jwt-key-change-in-production \
  marketplace-api
```

## Tests

```bash
go test ./...
```
