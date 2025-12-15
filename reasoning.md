# Reasoning & Key Decisions

## Goals
- Deliver a small, clean Go/Fiber API for user CRUD, storing `dob` and returning computed `age`.
- Keep layering clear (`handler` → `service` → `repository` → `sqlc`) and ops-friendly (Docker, env config).

## Architecture Choices
- **Fiber** for a fast, simple HTTP stack with minimal boilerplate.
- **Layering** to isolate concerns: handlers handle HTTP, services enforce validation/business logic (dob parsing, age computation), repositories wrap SQLC, models define DTOs and helpers.
- **SQLC** to keep SQL explicit while generating safe accessors. A stub of the generated code is included so the repo builds without running `sqlc generate`; real generation can replace it anytime.
- **Validation** via `go-playground/validator` to catch bad inputs at the edge.
- **Logging** with zap for structured, production-ready logs; middleware adds request IDs and duration metrics.
- **Age computation** done at read-time (`CalculateAge`) so age is always fresh and never stored.

## Error Handling & HTTP Semantics
- Fiber error handler centralizes JSON error responses.
- `404` for missing users, `400` for validation/parse errors, `201` for create, `204` for delete.

## Configuration
- `.env` via `godotenv`; defaults in code for sane local dev (`APP_PORT`, `DATABASE_URL`, `DB_MAX_CONNS`, `DB_MAX_IDLE`).

## Data & Migrations
- `db/migrations/000001_init.sql` defines the `users` table (id serial, name text, dob date).
- Queries in `db/sqlc/queries.sql`; `sqlc.yaml` points at migrations/queries and outputs to `db/sqlc`.

## Middleware
- `RequestID` to ensure traceability (`X-Request-ID`).
- `RequestLogger` to emit method/path/status/duration via zap.
- Fiber `recover` enabled to avoid panic crashes.

## Testing
- Unit test for age calculation (`internal/models/age_test.go`) covering edge cases (birthday passed, today, upcoming, future dob).

## Docker & Local DX
- `Dockerfile` (multi-stage, distroless runtime) and `docker-compose.yml` (API + Postgres).
- Keeps host setup light; Postgres defaults match `.env.example`.

## Trade-offs / Notes
- Included SQLC stub to keep the repo self-contained; teams running `sqlc generate` will overwrite it with real generated code.
- No `/` route by design; only `/users` endpoints. A `/healthz` can be added easily if desired.

