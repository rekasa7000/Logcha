# Logcha

A time tracking application designed for companies to monitor and manage rendered hours for On-the-Job Trainees (OJTs) and Interns.

## Features

- **Time Tracking**: Record time in and time out for daily attendance
- **Hours Management**: Track rendered hours and calculate remaining hours
- **Weekly Reports**: Generate total hours per week summaries
- **Progress Monitoring**: Monitor completion of required training hours

## Tech Stack

- **Backend**: Go with Fiber framework
- **Frontend**: React with TanStack

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Node.js 18 or higher
- npm or yarn

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/rekasa7000/Logcha.git
   cd Logcha
   ```

2. Backend setup:

   ```bash
   go mod tidy
   go run main.go
   ```

3. Frontend setup:
   ```bash
   cd client
   npm install
   npm run dev
   ```

## Project Structure

```
Logcha/
 client/           # React frontend with TanStack
 main.go           # Go backend with Fiber
 go.mod
 go.sum
 air.toml          # Air configuration for auto-restart
 README.md
```

## Contributing

This project is designed to help companies efficiently track and manage intern and OJT hours, ensuring compliance with training requirements and providing clear visibility into time commitments.
