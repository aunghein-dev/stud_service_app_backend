# Backend (Go + Chi + PostgreSQL)

## Quick start
1. Install dependencies:
   - Go 1.23+
   - `golang-migrate` CLI
2. Load environment from project root `.env`.
3. Run migrations:
   - `cd backend && make migrate-up DB_URL="$DB_URL"`
   - optional local demo data: `cd backend && make seed-dev-up DB_URL="$DB_URL"`
4. Start API:
   - `cd backend && make run`
5. Run tests:
   - `cd backend && make test`

## API docs
- Scalar API reference: `GET /docs`
- OpenAPI JSON: `GET /docs/openapi.json`
- Request and response schemas are generated from the Go DTO layer so payload docs stay close to the handlers.

## Core modules
- `auth / tenant workspaces`
- `students`
- `teachers`
- `class-courses`
- `optional fee items`
- `enrollments`
- `payments`
- `expenses`
- `receipts`
- `reports`
- `settings`

## Architecture
- `repository` handles SQL
- `service` handles business/transaction rules
- `handler` handles HTTP and DTO mapping
- `di` wires dependencies with Wire provider graph
- `middleware` now injects authenticated tenant context for all protected routes

## Migration strategy
- `make migrate-up` is production-safe and applies schema/index/auth changes only.
- Demo data is no longer part of the migration chain. It lives in `backend/seeds/development/up.sql`.
- `make seed-dev-up` uses `psql` to add a reusable sample tenant after migrations.
- `make seed-dev-down` removes the sample records without touching real production migrations.

## Receipt and transaction flow
- Enrollment create runs in a DB transaction
- Optional initial payment is saved in same transaction
- Receipt number allocated atomically from tenant-specific settings row
- Receipt payload is stored JSONB and can be re-opened/re-printed later

## Auth flow
- `POST /api/v1/auth/signup` creates or claims a tenant workspace and signs in the owner
- `POST /api/v1/auth/login` signs in an existing tenant user
- `GET /api/v1/auth/me` returns the current user and school workspace
- All other `/api/v1/*` endpoints require a bearer token
