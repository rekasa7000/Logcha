# Logcha Backend Setup

This is the Go Fiber backend for the Logcha time tracking application.

## Prerequisites

- Go 1.24.4 or higher
- PostgreSQL database (or Supabase account)

## Setup Instructions

### 1. Clone and Navigate to Project

```bash
cd Logcha
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Environment Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:

```env
# Server Configuration
PORT=4000
ENVIRONMENT=development

# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/logcha?sslmode=disable

# For Supabase:
# DATABASE_URL=postgres://postgres.your-project-ref:password@aws-0-region.pooler.supabase.com:5432/postgres

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

### 4. Database Setup

#### Option A: Local PostgreSQL

1. Install PostgreSQL
2. Create a database named `logcha`
3. Update the `DATABASE_URL` in your `.env` file

#### Option B: Supabase

1. Create a new project on [Supabase](https://supabase.com)
2. Go to Settings > Database
3. Copy the connection string and update `DATABASE_URL` in your `.env` file

### 5. Build and Run

```bash
# Build the application
go build -o bin/logcha cmd/server/main.go

# Run the application
./bin/logcha
```

Or run directly:

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:4000` (or the port specified in your `.env` file).

## API Endpoints

### Health Check
- `GET /api/v1/health` - Check if the API is running

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/me` - Get current user info (requires auth)

### Time Tracking
- `POST /api/v1/time/in` - Record time in (requires auth)
- `POST /api/v1/time/out` - Record time out (requires auth)
- `GET /api/v1/time/records/:traineeId` - Get time records for a trainee (requires auth)
- `GET /api/v1/time/today/:traineeId` - Get today's time record for a trainee (requires auth)

## Example API Usage

### Register a new user:
```bash
curl -X POST http://localhost:4000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "password": "password123",
    "role": "trainee"
  }'
```

### Login:
```bash
curl -X POST http://localhost:4000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### Time In (AM session):
```bash
curl -X POST http://localhost:4000/api/v1/time/in \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "trainee_id": 1,
    "session": "am"
  }'
```

## Project Structure

```
Logcha/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── handlers/
│   │   ├── auth_handler.go  # Authentication endpoints
│   │   └── time_record_handler.go # Time tracking endpoints
│   ├── middleware/
│   │   ├── auth.go          # JWT authentication middleware
│   │   ├── cors.go          # CORS middleware
│   │   └── logger.go        # Logging middleware
│   ├── models/
│   │   ├── user.go          # User model
│   │   ├── company.go       # Company model
│   │   ├── trainee.go       # Trainee model
│   │   ├── time_record.go   # Time record model
│   │   ├── weekly_summary.go # Weekly summary model
│   │   └── monthly_report.go # Monthly report model
│   ├── repository/
│   │   ├── interfaces.go    # Repository interfaces
│   │   ├── user_repository.go
│   │   ├── trainee_repository.go
│   │   └── time_record_repository.go
│   ├── services/
│   │   ├── auth_service.go  # Authentication business logic
│   │   └── time_record_service.go # Time tracking business logic
│   └── utils/
│       ├── jwt.go           # JWT utilities
│       ├── hash.go          # Password hashing
│       └── validation.go    # Request validation
├── pkg/
│   └── database/
│       └── postgres.go      # Database connection
├── .env.example             # Environment variables template
└── go.mod                   # Go module definition
```

## Features

### Implemented
- ✅ User authentication with JWT
- ✅ Time tracking (time in/out) for AM and PM sessions
- ✅ Database models for users, companies, trainees, and time records
- ✅ RESTful API endpoints
- ✅ CORS and logging middleware
- ✅ PostgreSQL/Supabase integration
- ✅ Auto-calculation of worked hours

### TODO (Future Enhancements)
- Weekly and monthly report generation
- Company admin features
- Trainee management endpoints
- Email notifications
- File upload for documents
- Advanced reporting and analytics

## Database Schema

The application follows the database schema outlined in `DATABASE_SCHEME.md`, supporting:

- Multi-company setup
- Different trainee types (paid interns, unpaid interns, OJT students)
- Flexible time tracking with AM/PM sessions
- Weekly hour limits and payment calculations
- Progress tracking for OJT students

## Troubleshooting

### Common Issues

1. **Database connection failed**: Check your `DATABASE_URL` in the `.env` file
2. **Port already in use**: Change the `PORT` in your `.env` file
3. **Build errors**: Run `go mod tidy` to ensure all dependencies are installed

### Development

For development with auto-restart, you can use Air:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with auto-restart
air
```

The `air.toml` configuration file is already included in the project.