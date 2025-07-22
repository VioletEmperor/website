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

### Content Storage Configuration
The application supports two storage modes for blog post content:

**Local Filesystem (Default)**:
- `STORAGE_MODE=local` (or unset)
- `POSTS_DIRECTORY=posts` (directory containing HTML files)

**Google Cloud Storage**:
- `STORAGE_MODE=gcs`
- `GCS_BUCKET_NAME=your-bucket-name` (required for GCS mode)
- `GCS_PREFIX=posts/` (optional prefix for post files)
- Requires `GOOGLE_APPLICATION_CREDENTIALS` environment variable or default GCP credentials

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

## Development Session Progress (2025-07-20)

### Completed Tasks
#### 1. Removed ADMIN_EMAILS Configuration
- **Files Modified**: `compose.yaml`, `internal/config/config.go`, `internal/handlers/handlers.go`, `static/js/admin-login.js`, `terraform/dev/run.tf`, `terraform/dev/secrets.tf`, `terraform/dev/variables.tf`
- **Changes**: Completely removed ADMIN_EMAILS environment variable and logic from all files
- **Result**: Firebase now handles all authentication - any authenticated user has admin access

#### 2. Enhanced Blog Posts Page
- **Files Created**: `static/css/posts.css`
- **Files Modified**: `templates/posts.html`, `templates/partials/post.html`
- **Features Added**:
  - Modern card-based layout matching website design
  - Displays title, author, creation date, last edited date, description
  - "Read Full Post" links to `/blog/post/{id}`
  - Responsive design with backdrop blur effects
  - Proper styling consistent with website theme

#### 3. Individual Post Viewing Setup
- **Files Modified**: `main.go`, `internal/handlers/handlers.go`
- **Route Added**: `GET /blog/post/{id}` -> `PostHandler`
- **Handler Created**: `PostHandler` extracts ID from URL path and calls `GetPost(id)`

### Pending Tasks
1. **Create individual post template** (`templates/post.html`) for full post display
2. **Verify GetPost method exists** in `internal/posts` repository
3. **Test blog functionality** end-to-end
4. **Fix any template/routing issues** that arise during testing

### Technical Notes
- Database schema includes: `id`, `title`, `author`, `created`, `edited`, `body`, `desc` fields
- Post cards use glassmorphism design with backdrop-filter blur
- Individual posts use path parameter extraction with `r.PathValue("id")`
- Template expects `Data{*post, "posts"}` structure for single post display