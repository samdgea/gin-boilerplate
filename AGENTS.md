# AGENTS.md

## Developer Commands
- Run server: `go run main.go`
- Live reload server: `./runApp` (requires `nodemon` installed globally)

## Environment & Prerequisites
- Copy `.env.example` to `.env` and fill all postgres/JWT values before running.
- Local PostgreSQL instance is required.

## Key Architecture Facts
- Database connection uses URL format in `db/db.go` to handle passwords with special characters.
- Auto-migration runs on startup via `models.MigrateModels(db)` called in `db.InitPostgres()`.
- Primary keys are UUIDs (`uuid.UUID`) generated in a Gorm hook `BeforeCreate` in `models/base.go`.
- Endpoint conventions under `/api/auth`:
  - `POST /api/auth/login` expects `userName` (camelCase) and `password` in body.
  - `POST /api/auth/refresh` expects `refresh_token` (snake_case) in body.
  - `POST /api/auth/logout` expects token in `Authorization` header.
- Token revocation is stateful: JWT access tokens map to a `TokenModel` database record. The middleware checks `TokenModel.IsActive` on every request.
