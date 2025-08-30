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

### Full Stack

- **Framework**: Next.js 15 with TypeScript and App Router
- **Database**: Supabase PostgreSQL with Row Level Security
- **Authentication**: Supabase Auth with JWT tokens
- **UI Library**: Shadcn/UI with Tailwind CSS 4.0
- **Form Handling**: React Hook Form with Zod validation
- **State Management**: React Context and Supabase client
- **Notifications**: Sonner for toast messages

### Development Tools

- **Package Manager**: Bun for fast dependency management
- **TypeScript**: Full type safety across the application
- **ESLint**: Code quality and consistency
- **Environment**: Development and production configurations

## Getting Started

### Prerequisites

- Node.js 18 or higher
- Bun (recommended) or npm
- Supabase account for database and authentication

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/rekasa7000/Logcha.git
   cd logcha
   ```

2. Install dependencies:

   ```bash
   bun install
   # or npm install
   ```

3. Set up environment variables:

   ```bash
   cp .env.example .env.local
   # Edit .env.local with your Supabase credentials
   ```

4. Set up Supabase:

   - Create a new Supabase project at [supabase.com](https://supabase.com)
   - Copy your project URL and anon key to `.env.local`
   - Run the SQL migration in `supabase-migration.sql` in your Supabase SQL editor
   - Run the admin seeder script: `bun run seed:admin`

5. Start the development server:

   ```bash
   bun dev
   # or npm run dev
   ```

### Environment Variables

Create a `.env.local` file in the root directory:

```env
# Supabase Configuration
NEXT_PUBLIC_SUPABASE_URL=your_supabase_project_url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key

# Optional: For server-side operations
SUPABASE_SERVICE_ROLE_KEY=your_supabase_service_role_key
```

## Project Structure

```
logcha/
├── app/                          # Next.js app directory
│   ├── auth/                     # Authentication pages
│   │   ├── login/page.tsx        # Login page
│   │   └── register/page.tsx     # Register page
│   ├── dashboard/                # Protected dashboard area
│   │   ├── companies/page.tsx    # Company management
│   │   ├── trainees/page.tsx     # Trainee management (planned)
│   │   ├── time-records/page.tsx # Time tracking (planned)
│   │   └── page.tsx             # Main dashboard
│   ├── layout.tsx               # Root layout with providers
│   ├── page.tsx                 # Home page with auth redirect
│   └── globals.css              # Global styles
├── components/                   # Reusable React components
│   ├── ui/                      # Shadcn UI components
│   ├── company-dialog.tsx       # Company form dialog
│   └── protected-route.tsx      # Route protection wrapper
├── lib/                         # Utility libraries
│   ├── auth-context.tsx         # Authentication context
│   ├── supabase.ts             # Supabase client configuration
│   ├── types.ts                # TypeScript type definitions
│   ├── validations.ts          # Zod validation schemas
│   └── utils.ts                # Utility functions
├── scripts/                     # Utility scripts
│   └── seed-admin.ts           # Admin user seeder
├── supabase-migration.sql       # Database schema migration
├── .env.example                # Environment variables template
├── .env.local                  # Local environment variables (gitignored)
├── components.json             # Shadcn UI configuration
├── package.json                # Node.js dependencies
├── tsconfig.json               # TypeScript configuration
├── tailwind.config.ts          # Tailwind CSS configuration
├── next.config.ts              # Next.js configuration
├── DATABASE_SCHEME.md          # Complete database schema
└── README.md                   # This file
```

## Application Features

The application provides comprehensive time tracking capabilities through a modern web interface:

### Authentication & Authorization

- **Supabase Auth**: Secure user authentication with email/password
- **Role-based Access**: Separate interfaces for company admins and trainees
- **Row Level Security**: Database-level security policies for data protection

### Company Management (Admin Only)

- Create and manage multiple companies
- Company profile with contact information and address
- Company-specific trainee management

### Trainee Management (Planned)

- Support for three trainee types:
  - **Paid Interns**: With hourly rates and weekly hour caps
  - **Unpaid Interns**: Hour tracking with weekly limits
  - **OJT Students**: Total required hours with progress tracking
- Individual trainee profiles with school and course information
- Flexible start/end dates and status management

### Time Tracking (Planned)

- **AM/PM Sessions**: Separate morning and afternoon time entries
- **Automatic Calculations**: Real-time hour calculations with generated columns
- **Validation**: Business rule enforcement (time logic, date restrictions)
- **Status Tracking**: Present, absent, half-day options

### Reporting & Analytics (Planned)

- **Weekly Summaries**: Automated weekly hour calculations with payroll
- **Monthly Reports**: Comprehensive DTR (Daily Time Record) reports
- **OJT Progress**: Visual progress tracking for required hours
- **Export Capabilities**: PDF and CSV export options

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

### Available Scripts

- `bun dev` - Start the development server with hot reload
- `bun build` - Build the application for production
- `bun start` - Start the production server
- `bun lint` - Run ESLint for code quality checks
- `bun run seed:admin` - Seed the database with an admin user

### Building for Production

```bash
# Build the application
bun build

# Start the production server
bun start
```

### Database Management

The application uses Supabase for database management:

- **Schema Migration**: Run `supabase-migration.sql` in your Supabase SQL editor
- **Row Level Security**: Policies are automatically applied for data protection
- **Real-time Updates**: Supabase provides real-time subscriptions for live data updates

### Development Workflow

1. Make your changes to the codebase
2. Test locally with `bun dev`
3. Run linting with `bun lint`
4. Test database changes in your Supabase project
5. Commit your changes with clear commit messages

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

- [Database Schema](DATABASE_SCHEME.md) - Complete database schema and business logic
- [Supabase Documentation](https://supabase.com/docs) - For database and authentication help
- [Next.js Documentation](https://nextjs.org/docs) - For framework-specific questions
- [Shadcn/UI Documentation](https://ui.shadcn.com) - For UI component usage
