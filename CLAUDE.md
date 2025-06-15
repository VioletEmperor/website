# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Local Development
```bash
# Start development environment with hot reload
docker compose up --build

# Run tests
go test ./...

# Build application
go build -o website
```

### Required Environment Variables
Before running, ensure these environment variables are set:
- `DATABASE_URL` - Full PostgreSQL connection string
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password  
- `DB_NAME` - Database name
- `PORT` - Server port (defaults to 8080 in Docker)

### Database Operations
```bash
# Reset database schema (via Docker)
docker compose down -v && docker compose up --build

# Connect to Cloud SQL database
./scripts/database-setup.sh <database_user>
```

### Deployment
```bash
# Deploy to Google Cloud Run
./scripts/update-cloud-run.sh
```

## Architecture Overview

This is a **server-side rendered Go web application** for a personal portfolio/blog with PostgreSQL backend.

### Project Structure
```
internal/
├── config/     - Environment configuration management
├── database/   - PostgreSQL connection handling (pgx/v5)
├── handlers/   - HTTP handlers with dependency injection
├── middleware/ - HTTP middleware (CORS, logging, stacking)
├── parse/      - HTML template parsing
└── posts/      - Blog post domain logic with repository pattern

templates/      - HTML templates (base layout + partials)
static/         - CSS, images, JavaScript assets
database/       - Schema and seed data
scripts/        - Cloud deployment scripts
```

### Key Patterns

**Dependency Injection**: Handlers are methods on `Env` struct containing database and template dependencies

**Repository Pattern**: Posts are accessed through `PostsRepository` interface with concrete PostgreSQL implementation

**Middleware Stack**: Uses custom middleware stacking in `internal/middleware/stack.go` - apply with `middleware.Stack(middleware1, middleware2)`

**Template System**: Pre-parsed HTML templates at startup with base layout + partials pattern

**Database**: Uses `pgx/v5/pgxpool` for connection pooling, schema recreated on container restart (no migrations)

### Development Workflow

1. **Adding Routes**: Create handler method on `Env` type, register in `main.go` with `mux.HandleFunc()`
2. **Database Changes**: Modify `database/init.sql`, restart containers to apply
3. **Templates**: Follow base → partials → page template hierarchy
4. **Middleware**: Add to stack in `main.go` using `middleware.Stack()`
5. **Configuration**: Add environment variables to `internal/config/config.go`

### Testing
- Tests run automatically during Docker build (`go test -v`)
- Use `pgxmock/v4` for database testing
- Config tests demonstrate environment variable mocking pattern