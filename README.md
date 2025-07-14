# Logcha

A comprehensive time tracking application designed for companies to monitor and manage rendered hours for On-the-Job Trainees (OJTs) and Interns. The system supports multiple trainee types including paid interns, unpaid interns, and OJT students with flexible hour tracking and reporting capabilities.

## Features

### Core Functionality
- **Multi-Session Time Tracking**: Record AM/PM sessions with automatic lunch break handling
- **Flexible Trainee Types**: Support for paid interns, unpaid interns, and OJT students
- **Weekly Hour Limits**: Company-defined maximum weekly hours for paid interns
- **OJT Progress Tracking**: Monitor completion of total required hours (e.g., 500 hours)
- **Automated Calculations**: Real-time hour calculations and payroll processing
- **Monthly Reports**: Generate comprehensive DTR (Daily Time Record) reports
- **Multi-Company Support**: Manage multiple companies and their trainees

### Business Logic
- **Paid Interns**: Hourly rate with weekly hour caps for payroll
- **Unpaid Interns**: Hour tracking with weekly limits, no payment
- **OJT Students**: Total hour requirements with progress tracking
- **Overtime Handling**: Track extra hours while respecting payment caps

## Tech Stack

### Backend
- **Framework**: Go 1.24.4 with Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT-based auth with bcrypt password hashing
- **Validation**: go-playground/validator for request validation
- **Configuration**: Environment-based config with godotenv

### Frontend
- **Framework**: Nuxt 3.17.6 (Vue 3.5.17)
- **UI Library**: Nuxt UI 3.2.0 with Tailwind CSS 4.0
- **Authentication**: Sidebase Nuxt Auth
- **State Management**: Composables-based approach
- **HTTP Client**: Custom API composable with JWT handling

### Development Tools
- **Hot Reload**: Air for Go backend auto-restart
- **Package Manager**: Bun for frontend dependencies
- **Environment**: Development and production configurations

## Getting Started

### Prerequisites

- Go 1.24.4 or higher
- Node.js 18 or higher
- PostgreSQL 15 or higher
- Bun (recommended) or npm

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/rekasa7000/Logcha.git
   cd Logcha
   ```

2. Backend setup:
   ```bash
   # Install Go dependencies
   go mod tidy
   
   # Set up environment variables
   cp .env.example .env
   # Edit .env with your database credentials
   
   # Run database migrations
   go run cmd/server/main.go
   ```

3. Frontend setup:
   ```bash
   cd client
   
   # Install dependencies
   bun install
   # or npm install
   
   # Start development server
   bun dev
   # or npm run dev
   ```

### Environment Variables

Create a `.env` file in the root directory:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=logcha

# JWT
JWT_SECRET=your_super_secret_jwt_key

# Server
PORT=8080
API_BASE_URL=http://localhost:8080

# Environment
ENVIRONMENT=development
```

## Project Structure

```
Logcha/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── handlers/
│   │   ├── auth_handler.go     # Authentication endpoints
│   │   └── time_record_handler.go  # Time tracking endpoints
│   ├── middleware/
│   │   ├── auth.go             # JWT authentication middleware
│   │   ├── cors.go             # CORS configuration
│   │   └── logger.go           # Request logging
│   ├── models/
│   │   ├── user.go             # User model
│   │   ├── company.go          # Company model
│   │   ├── trainee.go          # Trainee model
│   │   ├── time_record.go      # Time record model
│   │   ├── weekly_summary.go   # Weekly summary model
│   │   └── monthly_report.go   # Monthly report model
│   ├── repository/
│   │   ├── interfaces.go       # Repository interfaces
│   │   ├── user_repository.go  # User data access
│   │   ├── trainee_repository.go  # Trainee data access
│   │   └── time_record_repository.go  # Time record data access
│   ├── services/
│   │   ├── auth_service.go     # Authentication business logic
│   │   └── time_record_service.go  # Time tracking business logic
│   └── utils/
│       ├── hash.go             # Password hashing utilities
│       ├── jwt.go              # JWT token management
│       └── validation.go       # Input validation helpers
├── pkg/
│   └── database/
│       └── postgres.go         # Database connection setup
├── client/                     # Nuxt 3 frontend application
│   ├── assets/
│   │   └── css/
│   │       └── main.css        # Global styles
│   ├── composables/
│   │   ├── useApi.ts           # API client composable
│   │   └── useAuth.ts          # Authentication composable
│   ├── layouts/
│   │   └── default.vue         # Default layout component
│   ├── middleware/
│   │   ├── auth.ts             # Authentication middleware
│   │   └── guest.ts            # Guest-only middleware
│   ├── pages/
│   │   ├── index.vue           # Landing page
│   │   ├── login.vue           # Login page
│   │   ├── register.vue        # Registration page
│   │   └── dashboard.vue       # Main dashboard
│   ├── plugins/
│   │   └── api.client.ts       # API client plugin
│   ├── nuxt.config.ts          # Nuxt configuration
│   └── package.json            # Frontend dependencies
├── docs/
│   ├── DATABASE_SCHEME.md      # Complete database schema
│   ├── DEVELOPMENT_GUIDE.md    # Go & Fiber development guide
│   └── SETUP.md               # Detailed setup instructions
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── air.toml                   # Air configuration for hot reload
├── bun.lock                   # Bun lockfile
└── README.md                  # This file
```

## API Documentation

The backend provides a RESTful API with the following main endpoints:

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `GET /api/auth/me` - Get current user info

### Time Tracking
- `POST /api/time-records` - Create time record
- `GET /api/time-records` - Get time records
- `PUT /api/time-records/:id` - Update time record
- `DELETE /api/time-records/:id` - Delete time record

### Reports
- `GET /api/reports/weekly` - Weekly summary reports
- `GET /api/reports/monthly` - Monthly DTR reports
- `GET /api/reports/ojt-progress` - OJT progress tracking

## Database Schema

The application uses PostgreSQL with the following main tables:

- **users**: User authentication and profile information
- **companies**: Company/organization details
- **trainees**: Trainee profiles with type-specific configurations
- **time_records**: Daily time tracking with AM/PM sessions
- **weekly_summaries**: Calculated weekly hour summaries
- **monthly_reports**: Monthly DTR and payroll reports

For complete database schema details, see [DATABASE_SCHEME.md](DATABASE_SCHEME.md).

## Development

### Running with Hot Reload

Backend (with Air):
```bash
air
```

Frontend:
```bash
cd client
bun dev
```

### Building for Production

Backend:
```bash
go build -o bin/server cmd/server/main.go
```

Frontend:
```bash
cd client
bun build
```

### Testing

Run Go tests:
```bash
go test ./...
```

Run with coverage:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is designed to help companies efficiently track and manage intern and OJT hours, ensuring compliance with training requirements and providing clear visibility into time commitments.

## Support

For questions or support, please check the documentation:
- [Development Guide](DEVELOPMENT_GUIDE.md) - Complete Go & Fiber development guide
- [Database Schema](DATABASE_SCHEME.md) - Detailed database documentation
- [Setup Instructions](SETUP.md) - Detailed setup guide
