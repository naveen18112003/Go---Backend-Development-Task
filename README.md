# User Age API (Go + Fiber + SQLC)

REST API for managing users with dynamic age calculation. Built with GoFiber, Postgres, SQLC-style data layer, zap logging, and request validation.

## Implementation Plan
- Design: clear layering (`handler` -> `service` -> `repository` -> `sqlc`), zap logger, validator, middleware for request IDs and latency logging.
- Data: Postgres schema + SQLC queries; stubbed generated code included (`db/sqlc`) so the repo builds without running sqlc locally.
- Business logic: age computation via `time` package (`CalculateAge`) with unit tests.
- Transport: Fiber handlers with validation and clean HTTP codes; pagination on `/users`.
- Ops: Dockerfile + docker-compose for app + Postgres; `.env` sample in `config/env.example`.

## How to Run (Detailed)

### Option A: Docker (fastest)
Prereqs: Docker + docker-compose.
1) Build and start:
```sh
docker-compose up --build
```
2) Wait until logs show “database connected” and “starting server”.
3) API is available at `http://localhost:8080`.
4) Stop:
```sh
docker-compose down
```

### Option B: Local Go run (with your Postgres)
Prereqs: Go 1.21+, Postgres 15+, `psql` (or any client).
1) Create `.env` from template and edit as needed:
```sh
cp config/env.example .env
```
2) Ensure Postgres has a database (default uses `users`):
```sh
psql -U postgres -c "CREATE DATABASE users;"
```
3) Apply migration (simple path):
```sh
psql -U postgres -d users -f db/migrations/000001_init.sql
```
4) Verify `DATABASE_URL` in `.env` (defaults to `postgres://postgres:postgres@localhost:5432/users?sslmode=disable`).
5) Install deps and run:
```sh
go mod tidy
go run ./cmd/server
```
6) Server listens on `http://localhost:${APP_PORT}` (default 8080).

### Optional: Regenerate SQLC code
Prereqs: sqlc installed.
```sh
sqlc generate
```
Configuration: `sqlc.yaml`; queries: `db/sqlc/queries.sql`.

### Smoke-test the API
Create user:
```sh
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","dob":"1990-05-10"}'
```
Get by id (returns age):
```sh
curl http://localhost:8080/users/1
```
Update:
```sh
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","dob":"1991-03-15"}'
```
List (pagination):
```sh
curl "http://localhost:8080/users?limit=20&offset=0"
```
Delete:
```sh
curl -X DELETE http://localhost:8080/users/1
```

## API
- `POST /users` — create (`name`, `dob` as `YYYY-MM-DD`)
- `GET /users/:id` — fetch with `age`
- `PUT /users/:id` — update
- `DELETE /users/:id` — delete (204)
- `GET /users?limit=50&offset=0` — list with pagination + age

## Testing
```sh
go test ./internal/models
```

## Project Layout
See `/cmd/server` entrypoint, `internal/` for app layers, `db/` for migrations + SQLC artifacts, `config/` for env loading.

## Notes
- Zap logs include request duration and IDs via middleware.
- Age is always computed at read time using `time.Now()`.
- Default port `8080`; override via `APP_PORT`.

