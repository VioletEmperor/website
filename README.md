# Personal Portfolio & Blog Website

A modern, server-side rendered Go web application featuring a personal portfolio with an integrated blog system and admin dashboard.

## 🚀 Features

### Portfolio
- **About Section**: Personal information and background
- **Experience**: Professional work history
- **Education**: Academic background  
- **Skills**: Technical skills showcase with interactive carousel
- **Projects**: Portfolio of personal and professional projects
- **Contact Form**: Integrated contact form with email notifications

### Blog System
- **Public Blog**: Browse and read blog posts with modern card-based layout
- **Individual Post Views**: Full post display with HTML content support
- **Admin Dashboard**: Complete blog post management system
- **Content Storage**: Flexible storage with local filesystem or Google Cloud Storage support

### Admin Features
- **Firebase Authentication**: Secure login system
- **Post Management**: Create, edit, update, and delete blog posts
- **Content Upload**: Support for HTML file uploads
- **Dashboard Interface**: Modern admin interface for content management

## 🛠 Tech Stack

### Backend
- **Language**: Go 1.24
- **Database**: PostgreSQL with pgx/v5 connection pooling
- **Authentication**: Firebase Auth
- **Email**: Resend API for contact form notifications
- **Storage**: Local filesystem or Google Cloud Storage for blog content

### Frontend
- **Templates**: Server-side HTML templates with base layout + partials pattern
- **Styling**: Modern CSS with glassmorphism design effects
- **JavaScript**: Vanilla JS for interactive components

### Infrastructure
- **Containerization**: Docker & Docker Compose for local development
- **Cloud Deployment**: Google Cloud Run with Terraform infrastructure as code
- **Database**: Cloud SQL PostgreSQL for production
- **Monitoring**: OpenTelemetry integration

### Key Dependencies
```go
// Core dependencies
github.com/jackc/pgx/v5           // PostgreSQL driver
firebase.google.com/go/v4         // Firebase authentication
cloud.google.com/go/storage       // Google Cloud Storage
github.com/resend/resend-go/v2    // Email service
```

## 🏗 Architecture

### Project Structure
```
internal/
├── config/     - Environment configuration management
├── database/   - PostgreSQL connection handling
├── handlers/   - HTTP handlers with dependency injection
├── middleware/ - HTTP middleware (CORS, logging, auth)
├── parse/      - HTML template parsing
├── posts/      - Blog post domain logic with repository pattern
└── content/    - Content storage abstraction (filesystem/GCS)

templates/      - HTML templates (base layout + partials)
static/         - CSS, images, JavaScript assets
database/       - Schema and seed data
scripts/        - Deployment and utility scripts
terraform/      - Infrastructure as code (dev/prod environments)
```

### Key Patterns
- **Dependency Injection**: Handlers are methods on `Env` struct containing database and template dependencies
- **Repository Pattern**: Posts are accessed through `PostsRepository` interface with PostgreSQL implementation
- **Middleware Stack**: Custom middleware stacking with CORS, logging, and authentication
- **Content Abstraction**: Pluggable content storage supporting both local filesystem and Google Cloud Storage

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Go 1.24+ (for local development)
- Google Cloud SDK (for production deployment)

### Environment Variables
Create a `.env` file or set these environment variables:

```bash
# Database Configuration
DATABASE_URL=postgresql://user:password@localhost:5432/website
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_HOST=database

# Application Configuration
PORT=8080
PROJECT_ID=your-gcp-project-id
EMAIL_KEY=your-resend-api-key
FIREBASE_WEB_API_KEY=your-firebase-web-api-key

# Content Storage (choose one)
STORAGE_MODE=local              # or "gcs"
POSTS_DIRECTORY=posts           # for local mode
GCS_BUCKET_NAME=your-bucket     # for GCS mode
GCS_PREFIX=posts/               # optional, for GCS mode

# Authentication
GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
```

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd website
   ```

2. **Set up environment variables**
   ```bash
   # Copy and edit environment variables
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start the development environment**
   ```bash
   # Start with hot reload
   docker compose up --build
   ```

4. **Access the application**
   - Website: http://localhost
   - Admin Dashboard: http://localhost/admin
   - Database: localhost:5432

### Development Commands

```bash
# Run tests
go test ./...

# Build application locally
go build -o website

# Reset database (destroys all data)
docker compose down -v && docker compose up --build

# Connect to database
./scripts/database-setup.sh <database_user>
```

## ☁️ Cloud Deployment

### Google Cloud Platform Setup

1. **Prerequisites**
   - Google Cloud SDK installed and authenticated
   - Terraform installed
   - Required GCP APIs enabled:
     - Cloud Run API
     - Cloud SQL API
     - Container Registry API
     - Cloud Storage API (if using GCS storage)

2. **Infrastructure Deployment**
   ```bash
   # Navigate to terraform directory
   cd terraform/dev  # or terraform/prod

   # Initialize Terraform
   terraform init

   # Plan deployment
   terraform plan

   # Apply infrastructure
   terraform apply
   ```

3. **Application Deployment**
   ```bash
   # Deploy to Cloud Run
   ./scripts/update-cloud-run.sh
   ```

### Production Configuration

The application supports separate dev and prod environments with Terraform:

- **Development**: `terraform/dev/` - scaled-down resources for testing
- **Production**: `terraform/prod/` - production-ready configuration with high availability

Key production features:
- Auto-scaling Cloud Run service
- Cloud SQL with automated backups
- Container Registry for image storage
- Secret Manager for sensitive configuration
- Cloud Storage for blog content (optional)

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test ./internal/config
```

Testing includes:
- Unit tests for configuration management
- Database integration tests with pgxmock
- Handler tests for HTTP endpoints

## 📝 Development Workflow

1. **Adding Routes**: Create handler method on `Env` type, register in `main.go`
2. **Database Changes**: Modify `database/init.sql`, restart containers
3. **Templates**: Follow base → partials → page template hierarchy
4. **Middleware**: Add to stack in `main.go` using `middleware.Stack()`
5. **Configuration**: Add environment variables to `internal/config/config.go`

## 📚 API Endpoints

### Public Routes
- `GET /` - Homepage
- `GET /about` - About page
- `GET /blog/posts` - Blog posts listing
- `GET /blog/post/{id}` - Individual blog post
- `GET /contact` - Contact form
- `POST /contact` - Submit contact form

### Admin Routes (Authentication Required)
- `GET /admin/` - Admin homepage
- `GET /admin/login` - Login page
- `GET /admin/dashboard` - Admin dashboard
- `POST /admin/verify` - Verify Firebase token
- `GET /admin/posts` - List all posts
- `GET /admin/posts/{id}` - Get specific post
- `PUT /admin/posts/{id}` - Update post
- `DELETE /admin/posts/{id}` - Delete post
- `POST /admin/posts/upload` - Upload new post

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.