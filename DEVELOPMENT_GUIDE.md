# Complete Go & Fiber Backend Development Guide

## Table of Contents

1. [Go Fundamentals](#go-fundamentals)
2. [Fiber Framework Basics](#fiber-framework-basics)
3. [Project Structure](#project-structure)
4. [Database Integration](#database-integration)
5. [Authentication & Authorization](#authentication--authorization)
6. [Middleware](#middleware)
7. [API Development](#api-development)
8. [Cron Jobs & Background Tasks](#cron-jobs--background-tasks)
9. [File Operations](#file-operations)
10. [Testing](#testing)
11. [Deployment](#deployment)
12. [Best Practices](#best-practices)

---

## Go Fundamentals

### Installation & Setup

```bash
# Install Go
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Verify installation
go version

# Initialize a new project
mkdir my-fiber-app
cd my-fiber-app
go mod init my-fiber-app
```

### Basic Go Syntax

```go
package main

import (
    "fmt"
    "time"
)

// Variables
var globalVar string = "global"

// Constants
const PI = 3.14159

// Structs
type User struct {
    ID       int       `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Password string    `json:"-"` // Hidden from JSON
    CreatedAt time.Time `json:"created_at"`
}

// Methods
func (u *User) FullName() string {
    return fmt.Sprintf("User: %s", u.Name)
}

// Functions
func calculateAge(birthYear int) int {
    return time.Now().Year() - birthYear
}

// Error handling
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// Interfaces
type Database interface {
    Connect() error
    Close() error
}

// Goroutines and Channels
func processData(data chan string, done chan bool) {
    for item := range data {
        fmt.Println("Processing:", item)
        time.Sleep(1 * time.Second)
    }
    done <- true
}

func main() {
    // Basic types
    var name string = "John"
    age := 30
    isActive := true

    // Arrays and Slices
    var numbers [5]int = [5]int{1, 2, 3, 4, 5}
    slice := []string{"apple", "banana", "orange"}

    // Maps
    userMap := make(map[string]int)
    userMap["john"] = 25

    // Control structures
    if age >= 18 {
        fmt.Println("Adult")
    } else {
        fmt.Println("Minor")
    }

    // Loops
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }

    for index, value := range slice {
        fmt.Printf("Index: %d, Value: %s\n", index, value)
    }

    // Switch
    switch name {
    case "John":
        fmt.Println("Hello John!")
    case "Jane":
        fmt.Println("Hello Jane!")
    default:
        fmt.Println("Hello stranger!")
    }

    // Goroutines
    data := make(chan string, 3)
    done := make(chan bool)

    go processData(data, done)

    data <- "item1"
    data <- "item2"
    data <- "item3"
    close(data)

    <-done
}
```

---

## Fiber Framework Basics

### Installation

```bash
go get github.com/gofiber/fiber/v2
go get github.com/gofiber/fiber/v2/middleware/cors
go get github.com/gofiber/fiber/v2/middleware/logger
go get github.com/gofiber/fiber/v2/middleware/recover
```

### Basic Fiber App

```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    // Create fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
            })
        },
    })

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New())

    // Routes
    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Hello World!",
        })
    })

    // Start server
    log.Fatal(app.Listen(":3000"))
}
```

---

## Project Structure

```
my-fiber-app/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── user.go
│   │   └── product.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── models/
│   │   ├── user.go
│   │   └── product.go
│   ├── repository/
│   │   ├── interfaces.go
│   │   ├── user_repo.go
│   │   └── product_repo.go
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   └── product_service.go
│   └── utils/
│       ├── jwt.go
│       ├── hash.go
│       └── validation.go
├── pkg/
│   └── database/
│       └── postgres.go
├── migrations/
│   ├── 001_create_users_table.sql
│   └── 002_create_products_table.sql
├── scripts/
│   └── migrate.go
├── .env
├── .gitignore
├── go.mod
└── go.sum
```

### Configuration Management

```go
// internal/config/config.go
package config

import (
    "os"
    "strconv"
)

type Config struct {
    Port         string
    DatabaseURL  string
    JWTSecret    string
    Environment  string
    RedisURL     string
}

func Load() *Config {
    return &Config{
        Port:         getEnv("PORT", "3000"),
        DatabaseURL:  getEnv("DATABASE_URL", "postgres://user:password@localhost/dbname?sslmode=disable"),
        JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
        Environment:  getEnv("ENVIRONMENT", "development"),
        RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}
```

---

## Database Integration

### PostgreSQL with GORM

```bash
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-migrate/migrate/v4
```

```go
// pkg/database/postgres.go
package database

import (
    "fmt"
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type DB struct {
    *gorm.DB
}

func NewPostgresDB(dsn string) (*DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    // Get underlying sql.DB
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
    }

    // Set connection pool settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)

    return &DB{db}, nil
}

func (d *DB) AutoMigrate(models ...interface{}) error {
    return d.DB.AutoMigrate(models...)
}

func (d *DB) Close() error {
    sqlDB, err := d.DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}
```

### Models

```go
// internal/models/user.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"not null"`
    Email     string         `json:"email" gorm:"unique;not null"`
    Password  string         `json:"-" gorm:"not null"`
    Role      string         `json:"role" gorm:"default:user"`
    IsActive  bool           `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    Products []Product `json:"products,omitempty" gorm:"foreignKey:UserID"`
}

// internal/models/product.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    Name        string         `json:"name" gorm:"not null"`
    Description string         `json:"description"`
    Price       float64        `json:"price" gorm:"not null"`
    Stock       int            `json:"stock" gorm:"default:0"`
    UserID      uint           `json:"user_id" gorm:"not null"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    // Relationships
    User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
```

### Repository Pattern

```go
// internal/repository/interfaces.go
package repository

import (
    "context"
    "my-fiber-app/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uint) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, offset, limit int) ([]*models.User, error)
}

type ProductRepository interface {
    Create(ctx context.Context, product *models.Product) error
    GetByID(ctx context.Context, id uint) (*models.Product, error)
    GetByUserID(ctx context.Context, userID uint) ([]*models.Product, error)
    Update(ctx context.Context, product *models.Product) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, offset, limit int) ([]*models.Product, error)
}
```

```go
// internal/repository/user_repo.go
package repository

import (
    "context"
    "my-fiber-app/internal/models"
    "my-fiber-app/pkg/database"
)

type userRepository struct {
    db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
    var users []*models.User
    err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error
    return users, err
}
```

---

## Authentication & Authorization

### JWT Utilities

```go
// internal/utils/jwt.go
package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID uint, email, role, secret string) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Subject:   email,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, jwt.ErrInvalidKey
}
```

### Password Hashing

```go
// internal/utils/hash.go
package utils

import (
    "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### Authentication Service

```go
// internal/services/auth_service.go
package services

import (
    "context"
    "errors"
    "my-fiber-app/internal/models"
    "my-fiber-app/internal/repository"
    "my-fiber-app/internal/utils"
)

type AuthService struct {
    userRepo  repository.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
    return &AuthService{
        userRepo:  userRepo,
        jwtSecret: jwtSecret,
    }
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
    Name     string `json:"name" validate:"required,min=2"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type AuthResponse struct {
    Token string       `json:"token"`
    User  *models.User `json:"user"`
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
    user, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if !utils.CheckPassword(req.Password, user.Password) {
        return nil, errors.New("invalid credentials")
    }

    if !user.IsActive {
        return nil, errors.New("account is deactivated")
    }

    token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret)
    if err != nil {
        return nil, err
    }

    return &AuthResponse{
        Token: token,
        User:  user,
    }, nil
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
    // Check if user already exists
    if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
        return nil, errors.New("user already exists")
    }

    // Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    // Create user
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
        Role:     "user",
        IsActive: true,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // Generate token
    token, err := utils.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret)
    if err != nil {
        return nil, err
    }

    return &AuthResponse{
        Token: token,
        User:  user,
    }, nil
}
```

---

## Middleware

### Authentication Middleware

```go
// internal/middleware/auth.go
package middleware

import (
    "strings"
    "my-fiber-app/internal/utils"
    "github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtSecret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Missing authorization header",
            })
        }

        // Extract token from "Bearer <token>"
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token format",
            })
        }

        // Validate token
        claims, err := utils.ValidateToken(tokenString, jwtSecret)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }

        // Store user info in context
        c.Locals("user_id", claims.UserID)
        c.Locals("user_email", claims.Email)
        c.Locals("user_role", claims.Role)

        return c.Next()
    }
}

func RequireRole(role string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userRole := c.Locals("user_role")
        if userRole != role {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "error": "Insufficient permissions",
            })
        }
        return c.Next()
    }
}
```

### Rate Limiting

```go
// internal/middleware/rate_limit.go
package middleware

import (
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/limiter"
)

func RateLimitMiddleware() fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        100,
        Expiration: 1 * time.Minute,
        KeyGenerator: func(c *fiber.Ctx) string {
            return c.IP()
        },
        LimitReached: func(c *fiber.Ctx) error {
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error": "Rate limit exceeded",
            })
        },
    })
}
```

---

## API Development

### Handlers

```go
// internal/handlers/auth.go
package handlers

import (
    "my-fiber-app/internal/services"
    "my-fiber-app/internal/utils"
    "github.com/gofiber/fiber/v2"
    "github.com/go-playground/validator/v10"
)

type AuthHandler struct {
    authService *services.AuthService
    validator   *validator.Validate
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
        validator:   validator.New(),
    }
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req services.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if err := h.validator.Struct(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": utils.FormatValidationErrors(err),
        })
    }

    response, err := h.authService.Login(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(response)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req services.RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    if err := h.validator.Struct(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": utils.FormatValidationErrors(err),
        })
    }

    response, err := h.authService.Register(c.Context(), &req)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    userEmail := c.Locals("user_email").(string)
    userRole := c.Locals("user_role").(string)

    return c.JSON(fiber.Map{
        "id":    userID,
        "email": userEmail,
        "role":  userRole,
    })
}
```

### Routes

```go
// internal/routes/routes.go
package routes

import (
    "my-fiber-app/internal/handlers"
    "my-fiber-app/internal/middleware"
    "github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, jwtSecret string) {
    // API version
    api := app.Group("/api/v1")

    // Health check
    api.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "OK",
            "message": "Server is running",
        })
    })

    // Public routes
    auth := api.Group("/auth")
    auth.Post("/login", authHandler.Login)
    auth.Post("/register", authHandler.Register)

    // Protected routes
    protected := api.Group("/", middleware.AuthMiddleware(jwtSecret))
    protected.Get("/me", authHandler.Me)

    // User routes
    users := protected.Group("/users")
    users.Get("/", userHandler.GetUsers)
    users.Get("/:id", userHandler.GetUser)
    users.Put("/:id", userHandler.UpdateUser)
    users.Delete("/:id", userHandler.DeleteUser)

    // Admin routes
    admin := protected.Group("/admin", middleware.RequireRole("admin"))
    admin.Get("/users", userHandler.GetAllUsers)
    admin.Put("/users/:id/role", userHandler.UpdateUserRole)
}
```

---

## Cron Jobs & Background Tasks

### Cron Jobs

```bash
go get github.com/robfig/cron/v3
```

```go
// internal/jobs/scheduler.go
package jobs

import (
    "context"
    "log"
    "time"
    "my-fiber-app/internal/repository"
    "github.com/robfig/cron/v3"
)

type Scheduler struct {
    cron     *cron.Cron
    userRepo repository.UserRepository
}

func NewScheduler(userRepo repository.UserRepository) *Scheduler {
    return &Scheduler{
        cron:     cron.New(),
        userRepo: userRepo,
    }
}

func (s *Scheduler) Start() {
    // Run every day at midnight
    s.cron.AddFunc("0 0 * * *", s.cleanupInactiveUsers)

    // Run every hour
    s.cron.AddFunc("0 * * * *", s.generateReports)

    // Run every 5 minutes
    s.cron.AddFunc("*/5 * * * *", s.healthCheck)

    s.cron.Start()
    log.Println("Cron jobs started")
}

func (s *Scheduler) Stop() {
    s.cron.Stop()
    log.Println("Cron jobs stopped")
}

func (s *Scheduler) cleanupInactiveUsers() {
    ctx := context.Background()
    log.Println("Running cleanup inactive users job...")

    // Implement your cleanup logic here
    // For example, soft delete users inactive for 30 days

    log.Println("Cleanup job completed")
}

func (s *Scheduler) generateReports() {
    ctx := context.Background()
    log.Println("Generating reports...")

    // Implement report generation logic

    log.Println("Reports generated")
}

func (s *Scheduler) healthCheck() {
    log.Println("Health check completed")
}
```

### Background Workers

```go
// internal/workers/worker.go
package workers

import (
    "context"
    "log"
    "time"
)

type Worker struct {
    jobQueue   chan Job
    workerPool chan chan Job
    quit       chan bool
}

type Job struct {
    Type string
    Data interface{}
}

func NewWorker(workerPool chan chan Job) *Worker {
    return &Worker{
        jobQueue:   make(chan Job),
        workerPool: workerPool,
        quit:       make(chan bool),
    }
}

func (w *Worker) Start() {
    go func() {
        for {
            w.workerPool <- w.jobQueue

            select {
            case job := <-w.jobQueue:
                w.processJob(job)
            case <-w.quit:
                return
            }
        }
    }()
}

func (w *Worker) Stop() {
    w.quit <- true
}

func (w *Worker) processJob(job Job) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    switch job.Type {
    case "send_email":
        w.sendEmail(ctx, job.Data)
    case "process_image":
        w.processImage(ctx, job.Data)
    default:
        log.Printf("Unknown job type: %s", job.Type)
    }
}

func (w *Worker) sendEmail(ctx context.Context, data interface{}) {
    log.Printf("Sending email: %v", data)
    // Implement email sending logic
    time.Sleep(2 * time.Second) // Simulate work
}

func (w *Worker) processImage(ctx context.Context, data interface{}) {
    log.Printf("Processing image: %v", data)
    // Implement image processing logic
    time.Sleep(5 * time.Second) // Simulate work
}

// Worker Pool
type WorkerPool struct {
    workers    []*Worker
    jobQueue   chan Job
    workerPool chan chan Job
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
    return &WorkerPool{
        workers:    make([]*Worker, maxWorkers),
        jobQueue:   make(chan Job, 100),
        workerPool: make(chan chan Job, maxWorkers),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < len(wp.workers); i++ {
        wp.workers[i] = NewWorker(wp.workerPool)
        wp.workers[i].Start()
    }

    go wp.dispatch()
}

func (wp *WorkerPool) dispatch() {
    for {
        select {
        case job := <-wp.jobQueue:
            go func(job Job) {
                jobQueue := <-wp.workerPool
                jobQueue <- job
            }(job)
        }
    }
}

func (wp *WorkerPool) AddJob(job Job) {
    wp.jobQueue <- job
}

func (wp *WorkerPool) Stop() {
    for _, worker := range wp.workers {
        worker.Stop()
    }
}
```

---

## File Operations

### File Upload Handler

```bash
go get github.com/gofiber/fiber/v2/middleware/filesystem
```

```go
// internal/handlers/file.go
package handlers

import (
    "fmt"
    "io"
    "mime/multipart"
    "os"
    "path/filepath"
    "strings"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

type FileHandler struct {
    uploadDir string
    maxSize   int64
}

func NewFileHandler(uploadDir string, maxSize int64) *FileHandler {
    // Create upload directory if it doesn't exist
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        panic(err)
    }

    return &FileHandler{
        uploadDir: uploadDir,
        maxSize:   maxSize,
    }
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }

    // Check file size
    if file.Size > h.maxSize {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": fmt.Sprintf("File too large. Max size: %d bytes", h.maxSize),
        })
    }

    // Validate file type
    if !h.isValidFileType(file) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid file type",
        })
    }

    // Generate unique filename
    filename := h.generateFilename(file.Filename)
    filepath := filepath.Join(h.uploadDir, filename)

    // Save file
    if err := c.SaveFile(file, filepath); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to save file",
        })
    }

    return c.JSON(fiber.Map{
        "message":  "File uploaded successfully",
        "filename": filename,
        "size":     file.Size,
        "path":     filepath,
    })
}

func (h *FileHandler) UploadMultipleFiles(c *fiber.Ctx) error {
    form, err := c.MultipartForm()
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to parse multipart form",
        })
    }

    files := form.File["files"]
    if len(files) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "No files uploaded",
        })
    }

    var uploadedFiles []map[string]interface{}
    var errors []string

    for _, file := range files {
        if file.Size > h.maxSize {
            errors = append(errors, fmt.Sprintf("File %s too large", file.Filename))
            continue
        }

        if !h.isValidFileType(file) {
            errors = append(errors, fmt.Sprintf("Invalid file type for %s", file.Filename))
            continue
        }

        filename := h.generateFilename(file.Filename)
        filepath := filepath.Join(h.uploadDir, filename)

        if err := c.SaveFile(file, filepath); err != nil {
            errors = append(errors, fmt.Sprintf("Failed to save %s", file.Filename))
            continue
        }

        uploadedFiles = append(uploadedFiles, map[string]interface{}{
            "original_name": file.Filename,
            "filename":      filename,
            "size":          file.Size,
            "path":          filepath,
        })
    }

    response := fiber.Map{
        "uploaded_files": uploadedFiles,
        "total_uploaded": len(uploadedFiles),
    }

    if len(errors) > 0 {
        response["errors"] = errors
    }

    return c.JSON(response)
}

func (h *FileHandler) DownloadFile(c *fiber.Ctx) error {
    filename := c.Params("filename")
    if filename == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Filename required",
        })
    }

    filepath := filepath.Join(h.uploadDir, filename)

    // Check if file exists
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "File not found",
        })
    }

    return c.SendFile(filepath)
}

func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
    filename := c.Params("filename")
    if filename == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Filename required",
        })
    }

    filepath := filepath.Join(h.uploadDir, filename)

    if err := os.Remove(filepath); err != nil {
        if os.IsNotExist(err) {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "File not found",
            })
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to delete file",
        })
    }

    return c.JSON(fiber.Map{
        "message": "File deleted successfully",
    })
}

func (h *FileHandler) isValidFileType(file *multipart.FileHeader) bool {
    allowedTypes := map[string]bool{
        "image/jpeg":      true,
        "image/jpg":       true,
        "image/png":       true,
        "image/gif":       true,
        "application/pdf": true,
        "text/plain":      true,
        "text/csv":        true,
    }

    // Get file extension
    ext := strings.ToLower(filepath.Ext(file.Filename))
    allowedExts := map[string]bool{
        ".jpg":  true,
        ".jpeg": true,
        ".png":  true,
        ".gif":  true,
        ".pdf":  true,
        ".txt":  true,
        ".csv":  true,
    }

    return allowedExts[ext]
}

func (h *FileHandler) generateFilename(originalName string) string {
    ext := filepath.Ext(originalName)
    name := strings.TrimSuffix(originalName, ext)
    timestamp := time.Now().Unix()
    uuid := uuid.New().String()[:8]

    return fmt.Sprintf("%s_%d_%s%s", name, timestamp, uuid, ext)
}
```

### Image Processing

```bash
go get github.com/disintegration/imaging
```

```go
// internal/utils/image.go
package utils

import (
    "fmt"
    "image"
    "path/filepath"
    "github.com/disintegration/imaging"
)

type ImageProcessor struct {
    outputDir string
}

func NewImageProcessor(outputDir string) *ImageProcessor {
    return &ImageProcessor{outputDir: outputDir}
}

func (p *ImageProcessor) ResizeImage(inputPath string, width, height int) (string, error) {
    // Open the image
    img, err := imaging.Open(inputPath)
    if err != nil {
        return "", fmt.Errorf("failed to open image: %w", err)
    }

    // Resize the image
    resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

    // Generate output filename
    filename := fmt.Sprintf("resized_%dx%d_%s", width, height, filepath.Base(inputPath))
    outputPath := filepath.Join(p.outputDir, filename)

    // Save the resized image
    if err := imaging.Save(resizedImg, outputPath); err != nil {
        return "", fmt.Errorf("failed to save resized image: %w", err)
    }

    return outputPath, nil
}

func (p *ImageProcessor) CreateThumbnail(inputPath string, size int) (string, error) {
    img, err := imaging.Open(inputPath)
    if err != nil {
        return "", fmt.Errorf("failed to open image: %w", err)
    }

    // Create thumbnail (square crop)
    thumbnail := imaging.Fill(img, size, size, imaging.Center, imaging.Lanczos)

    filename := fmt.Sprintf("thumb_%dx%d_%s", size, size, filepath.Base(inputPath))
    outputPath := filepath.Join(p.outputDir, filename)

    if err := imaging.Save(thumbnail, outputPath); err != nil {
        return "", fmt.Errorf("failed to save thumbnail: %w", err)
    }

    return outputPath, nil
}

func (p *ImageProcessor) ConvertFormat(inputPath, outputFormat string) (string, error) {
    img, err := imaging.Open(inputPath)
    if err != nil {
        return "", fmt.Errorf("failed to open image: %w", err)
    }

    // Generate output filename with new extension
    baseFilename := filepath.Base(inputPath)
    name := baseFilename[:len(baseFilename)-len(filepath.Ext(baseFilename))]
    filename := fmt.Sprintf("%s.%s", name, outputFormat)
    outputPath := filepath.Join(p.outputDir, filename)

    // Save with new format
    if err := imaging.Save(img, outputPath); err != nil {
        return "", fmt.Errorf("failed to save converted image: %w", err)
    }

    return outputPath, nil
}
```

### CSV Processing

```go
// internal/utils/csv.go
package utils

import (
    "encoding/csv"
    "fmt"
    "os"
    "reflect"
    "strconv"
)

type CSVProcessor struct{}

func NewCSVProcessor() *CSVProcessor {
    return &CSVProcessor{}
}

func (p *CSVProcessor) ReadCSV(filepath string) ([][]string, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, fmt.Errorf("failed to open CSV file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("failed to read CSV: %w", err)
    }

    return records, nil
}

func (p *CSVProcessor) WriteCSV(filepath string, data [][]string) error {
    file, err := os.Create(filepath)
    if err != nil {
        return fmt.Errorf("failed to create CSV file: %w", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, record := range data {
        if err := writer.Write(record); err != nil {
            return fmt.Errorf("failed to write CSV record: %w", err)
        }
    }

    return nil
}

func (p *CSVProcessor) StructToCSV(data interface{}, filepath string) error {
    v := reflect.ValueOf(data)
    if v.Kind() != reflect.Slice {
        return fmt.Errorf("data must be a slice")
    }

    if v.Len() == 0 {
        return fmt.Errorf("data slice is empty")
    }

    // Get the type of the first element
    elemType := v.Index(0).Type()
    var headers []string
    var records [][]string

    // Extract headers from struct fields
    for i := 0; i < elemType.NumField(); i++ {
        field := elemType.Field(i)
        headers = append(headers, field.Name)
    }
    records = append(records, headers)

    // Extract data
    for i := 0; i < v.Len(); i++ {
        elem := v.Index(i)
        var record []string
        for j := 0; j < elem.NumField(); j++ {
            field := elem.Field(j)
            record = append(record, p.fieldToString(field))
        }
        records = append(records, record)
    }

    return p.WriteCSV(filepath, records)
}

func (p *CSVProcessor) fieldToString(field reflect.Value) string {
    switch field.Kind() {
    case reflect.String:
        return field.String()
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return strconv.FormatInt(field.Int(), 10)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return strconv.FormatUint(field.Uint(), 10)
    case reflect.Float32, reflect.Float64:
        return strconv.FormatFloat(field.Float(), 'f', -1, 64)
    case reflect.Bool:
        return strconv.FormatBool(field.Bool())
    default:
        return fmt.Sprintf("%v", field.Interface())
    }
}
```

---

## Testing

### Unit Tests

```go
// internal/services/auth_service_test.go
package services

import (
    "context"
    "testing"
    "errors"
    "my-fiber-app/internal/models"
    "my-fiber-app/internal/utils"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    args := m.Called(ctx, email)
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
    args := m.Called(ctx, offset, limit)
    return args.Get(0).([]*models.User), args.Error(1)
}

func TestAuthService_Login(t *testing.T) {
    mockRepo := new(MockUserRepository)
    authService := NewAuthService(mockRepo, "test-secret")

    t.Run("successful login", func(t *testing.T) {
        hashedPassword, _ := utils.HashPassword("password123")
        user := &models.User{
            ID:       1,
            Email:    "test@example.com",
            Password: hashedPassword,
            Role:     "user",
            IsActive: true,
        }

        mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

        req := &LoginRequest{
            Email:    "test@example.com",
            Password: "password123",
        }

        response, err := authService.Login(context.Background(), req)

        assert.NoError(t, err)
        assert.NotNil(t, response)
        assert.NotEmpty(t, response.Token)
        assert.Equal(t, user.Email, response.User.Email)
        mockRepo.AssertExpectations(t)
    })

    t.Run("invalid credentials", func(t *testing.T) {
        mockRepo.On("GetByEmail", mock.Anything, "wrong@example.com").Return((*models.User)(nil), errors.New("user not found"))

        req := &LoginRequest{
            Email:    "wrong@example.com",
            Password: "password123",
        }

        response, err := authService.Login(context.Background(), req)

        assert.Error(t, err)
        assert.Nil(t, response)
        assert.Contains(t, err.Error(), "invalid credentials")
        mockRepo.AssertExpectations(t)
    })

    t.Run("inactive user", func(t *testing.T) {
        hashedPassword, _ := utils.HashPassword("password123")
        user := &models.User{
            ID:       1,
            Email:    "test@example.com",
            Password: hashedPassword,
            Role:     "user",
            IsActive: false, // Inactive user
        }

        mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil)

        req := &LoginRequest{
            Email:    "test@example.com",
            Password: "password123",
        }

        response, err := authService.Login(context.Background(), req)

        assert.Error(t, err)
        assert.Nil(t, response)
        assert.Contains(t, err.Error(), "account is deactivated")
        mockRepo.AssertExpectations(t)
    })
}

func TestAuthService_Register(t *testing.T) {
    mockRepo := new(MockUserRepository)
    authService := NewAuthService(mockRepo, "test-secret")

    t.Run("successful registration", func(t *testing.T) {
        mockRepo.On("GetByEmail", mock.Anything, "new@example.com").Return((*models.User)(nil), errors.New("user not found"))
        mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

        req := &RegisterRequest{
            Name:     "John Doe",
            Email:    "new@example.com",
            Password: "password123",
        }

        response, err := authService.Register(context.Background(), req)

        assert.NoError(t, err)
        assert.NotNil(t, response)
        assert.NotEmpty(t, response.Token)
        assert.Equal(t, req.Email, response.User.Email)
        mockRepo.AssertExpectations(t)
    })

    t.Run("user already exists", func(t *testing.T) {
        existingUser := &models.User{
            ID:    1,
            Email: "existing@example.com",
        }

        mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)

        req := &RegisterRequest{
            Name:     "John Doe",
            Email:    "existing@example.com",
            Password: "password123",
        }

        response, err := authService.Register(context.Background(), req)

        assert.Error(t, err)
        assert.Nil(t, response)
        assert.Contains(t, err.Error(), "user already exists")
        mockRepo.AssertExpectations(t)
    })
}
```

### Integration Tests

```go
// tests/integration/auth_test.go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "my-fiber-app/internal/config"
    "my-fiber-app/internal/handlers"
    "my-fiber-app/internal/services"
    "my-fiber-app/pkg/database"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
    suite.Suite
    app        *fiber.App
    db         *database.DB
    authHandler *handlers.AuthHandler
}

func (suite *AuthTestSuite) SetupSuite() {
    // Setup test database
    cfg := config.Load()
    db, err := database.NewPostgresDB(cfg.DatabaseURL)
    suite.Require().NoError(err)

    suite.db = db

    // Setup services and handlers
    userRepo := repository.NewUserRepository(db)
    authService := services.NewAuthService(userRepo, cfg.JWTSecret)
    suite.authHandler = handlers.NewAuthHandler(authService)

    // Setup Fiber app
    suite.app = fiber.New()
    suite.app.Post("/auth/register", suite.authHandler.Register)
    suite.app.Post("/auth/login", suite.authHandler.Login)
}

func (suite *AuthTestSuite) TearDownSuite() {
    suite.db.Close()
}

func (suite *AuthTestSuite) TestRegisterAndLogin() {
    // Test registration
    registerData := map[string]string{
        "name":     "Test User",
        "email":    "test@example.com",
        "password": "password123",
    }

    registerJSON, _ := json.Marshal(registerData)
    req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(registerJSON))
    req.Header.Set("Content-Type", "application/json")

    resp, _ := suite.app.Test(req)
    assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

    var registerResponse map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&registerResponse)
    assert.NotEmpty(suite.T(), registerResponse["token"])

    // Test login
    loginData := map[string]string{
        "email":    "test@example.com",
        "password": "password123",
    }

    loginJSON, _ := json.Marshal(loginData)
    req = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginJSON))
    req.Header.Set("Content-Type", "application/json")

    resp, _ = suite.app.Test(req)
    assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

    var loginResponse map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&loginResponse)
    assert.NotEmpty(suite.T(), loginResponse["token"])
}

func TestAuthTestSuite(t *testing.T) {
    suite.Run(t, new(AuthTestSuite))
}
```

### Test Utilities

```go
// tests/utils/test_utils.go
package utils

import (
    "my-fiber-app/internal/models"
    "my-fiber-app/internal/utils"
    "time"
)

func CreateTestUser(email, password string) *models.User {
    hashedPassword, _ := utils.HashPassword(password)
    return &models.User{
        ID:        1,
        Name:      "Test User",
        Email:     email,
        Password:  hashedPassword,
        Role:      "user",
        IsActive:  true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func CreateTestProduct(name string, price float64, userID uint) *models.Product {
    return &models.Product{
        ID:          1,
        Name:        name,
        Description: "Test product description",
        Price:       price,
        Stock:       10,
        UserID:      userID,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
}
```

---

## Deployment

### Docker Configuration

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Create uploads directory
RUN mkdir -p uploads

EXPOSE 3000

CMD ["./main"]
```

```yaml
# docker-compose.yml
version: "3.8"

services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DATABASE_URL=postgres://user:password@db:5432/myapp?sslmode=disable
      - JWT_SECRET=your-super-secret-jwt-key
      - ENVIRONMENT=production
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
    volumes:
      - ./uploads:/root/uploads
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=myapp
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```

### Nginx Configuration

```nginx
# nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream app {
        server app:3000;
    }

    server {
        listen 80;
        server_name yourdomain.com;

        location / {
            proxy_pass http://app;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # File upload size limit
        client_max_body_size 10M;

        # Static files
        location /uploads/ {
            alias /root/uploads/;
            expires 30d;
            add_header Cache-Control "public, immutable";
        }
    }
}
```

### Environment Variables

```bash
# .env
PORT=3000
DATABASE_URL=postgres://user:password@localhost:5432/myapp?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
ENVIRONMENT=development
REDIS_URL=redis://localhost:6379
UPLOAD_DIR=./uploads
MAX_FILE_SIZE=10485760
```

### GitHub Actions CI/CD

```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: password
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.21

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: go test -v ./...
        env:
          DATABASE_URL: postgres://postgres:password@localhost:5432/testdb?sslmode=disable

      - name: Run linter
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest

  build-and-deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
      - uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          push: true
          tags: yourusername/myapp:latest

      - name: Deploy to production
        run: |
          # Add your deployment script here
          echo "Deploying to production..."
```

---

## Best Practices

### Code Organization

1. **Follow Clean Architecture**: Separate concerns into layers (handlers, services, repositories)
2. **Use Dependency Injection**: Make components testable and maintainable
3. **Implement Interfaces**: Define contracts for better abstraction
4. **Handle Errors Properly**: Use proper error handling and logging
5. **Validate Input**: Always validate and sanitize user input

### Security Best Practices

```go
// internal/middleware/security.go
package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/helmet"
    "github.com/gofiber/fiber/v2/middleware/cors"
)

func SecurityMiddleware() fiber.Handler {
    return helmet.New(helmet.Config{
        XSSProtection:         "1; mode=block",
        ContentTypeNosniff:    "nosniff",
        XFrameOptions:         "DENY",
        ReferrerPolicy:        "no-referrer",
        CrossOriginEmbedderPolicy: "require-corp",
        CrossOriginOpenerPolicy:   "same-origin",
        CrossOriginResourcePolicy: "same-origin",
        OriginAgentCluster:        "?1",
        XDNSPrefetchControl:       "off",
        XDownloadOptions:          "noopen",
        XPermittedCrossDomain:     "none",
    })
}

func CORSMiddleware() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins:     "https://yourdomain.com",
        AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders:     "Origin,Content-Type,Authorization",
        AllowCredentials: true,
        MaxAge:           86400,
    })
}
```

### Performance Optimization

```go
// internal/middleware/cache.go
package middleware

import (
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cache"
)

func CacheMiddleware() fiber.Handler {
    return cache.New(cache.Config{
        Next: func(c *fiber.Ctx) bool {
            // Skip cache for POST, PUT, DELETE requests
            return c.Method() != "GET"
        },
        Expiration:   30 * time.Minute,
        CacheHeader:  "X-Cache",
        CacheControl: true,
    })
}
```

### Database Optimization

```go
// internal/repository/optimized_queries.go
package repository

import (
    "context"
    "my-fiber-app/internal/models"
    "my-fiber-app/pkg/database"
)

type OptimizedQueries struct {
    db *database.DB
}

func NewOptimizedQueries(db *database.DB) *OptimizedQueries {
    return &OptimizedQueries{db: db}
}

// Use indexes and proper query optimization
func (q *OptimizedQueries) GetUserWithProducts(ctx context.Context, userID uint) (*models.User, error) {
    var user models.User
    err := q.db.WithContext(ctx).
        Preload("Products").
        Where("id = ? AND is_active = ?", userID, true).
        First(&user).Error

    return &user, err
}

// Use batch operations for better performance
func (q *OptimizedQueries) CreateProductsBatch(ctx context.Context, products []*models.Product) error {
    return q.db.WithContext(ctx).CreateInBatches(products, 100).Error
}

// Use raw queries for complex operations
func (q *OptimizedQueries) GetTopSellingProducts(ctx context.Context, limit int) ([]*models.Product, error) {
    var products []*models.Product
    err := q.db.WithContext(ctx).
        Raw(`
            SELECT p.*, COUNT(o.id) as order_count
            FROM products p
            LEFT JOIN orders o ON p.id = o.product_id
            GROUP BY p.id
            ORDER BY order_count DESC
            LIMIT ?
        `, limit).
        Scan(&products).Error

    return products, err
}
```

### Logging and Monitoring

```go
// internal/utils/logger.go
package utils

import (
    "os"
    "github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
    Logger = logrus.New()

    // Set log level
    level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
    if err != nil {
        level = logrus.InfoLevel
    }
    Logger.SetLevel(level)

    // Set formatter
    Logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
    })

    // Set output
    Logger.SetOutput(os.Stdout)
}

// Middleware for request logging
func LoggerMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()

        // Process request
        err := c.Next()

        // Log request
        Logger.WithFields(logrus.Fields{
            "method":     c.Method(),
            "path":       c.Path(),
            "status":     c.Response().StatusCode(),
            "duration":   time.Since(start).String(),
            "ip":         c.IP(),
            "user_agent": c.Get("User-Agent"),
        }).Info("Request processed")

        return err
    }
}
```

### Error Handling

```go
// internal/utils/errors.go
package utils

import (
    "errors"
    "github.com/gofiber/fiber/v2"
)

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
    return e.Message
}

func NewAPIError(code int, message string) *APIError {
    return &APIError{
        Code:    code,
        Message: message,
    }
}

func ErrorHandler(c *fiber.Ctx, err error) error {
    var apiErr *APIError

    if errors.As(err, &apiErr) {
        return c.Status(apiErr.Code).JSON(apiErr)
    }

    // Handle Fiber errors
    if e, ok := err.(*fiber.Error); ok {
        return c.Status(e.Code).JSON(fiber.Map{
            "error": e.Message,
        })
    }

    // Default error
    Logger.Error("Unhandled error: ", err)
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": "Internal server error",
    })
}
```

### Configuration Management

```go
// internal/config/advanced_config.go
package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
    "github.com/joho/godotenv"
)

type DatabaseConfig struct {
    Host            string
    Port            int
    User            string
    Password        string
    Database        string
    SSLMode         string
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

type ServerConfig struct {
    Port         string
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}

type AppConfig struct {
    Environment string
    Debug       bool
    JWTSecret   string
    Server      ServerConfig
    Database    DatabaseConfig
    Redis       RedisConfig
    UploadDir   string
    MaxFileSize int64
}

func LoadConfig() (*AppConfig, error) {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        // .env file is optional
    }

    config := &AppConfig{
        Environment: getEnv("ENVIRONMENT", "development"),
        Debug:       getEnvAsBool("DEBUG", false),
        JWTSecret:   getEnv("JWT_SECRET", ""),
        Server: ServerConfig{
            Port:         getEnv("PORT", "3000"),
            ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", 10*time.Second),
            WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", 10*time.Second),
            IdleTimeout:  getEnvAsDuration("IDLE_TIMEOUT", 60*time.Second),
        },
        Database: DatabaseConfig{
            Host:            getEnv("DB_HOST", "localhost"),
            Port:            getEnvAsInt("DB_PORT", 5432),
            User:            getEnv("DB_USER", ""),
            Password:        getEnv("DB_PASSWORD", ""),
            Database:        getEnv("DB_NAME", ""),
            SSLMode:         getEnv("DB_SSLMODE", "disable"),
            MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
            MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
            ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnvAsInt("REDIS_PORT", 6379),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       getEnvAsInt("REDIS_DB", 0),
        },
        UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
        MaxFileSize: getEnvAsInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB
    }

    if config.JWTSecret == "" {
        return nil, fmt.Errorf("JWT_SECRET is required")
    }

    return config, nil
}

func getEnvAsBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
```

### Complete Main Application

```go
// cmd/server/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "my-fiber-app/internal/config"
    "my-fiber-app/internal/handlers"
    "my-fiber-app/internal/jobs"
    "my-fiber-app/internal/middleware"
    "my-fiber-app/internal/models"
    "my-fiber-app/internal/repository"
    "my-fiber-app/internal/routes"
    "my-fiber-app/internal/services"
    "my-fiber-app/internal/utils"
    "my-fiber-app/internal/workers"
    "my-fiber-app/pkg/database"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
    // Initialize logger
    utils.InitLogger()

    // Load configuration
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatal("Failed to load configuration:", err)
    }

    // Connect to database
    db, err := database.NewPostgresDB(cfg.Database)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Auto-migrate database
    if err := db.AutoMigrate(&models.User{}, &models.Product{}); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }

    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    productRepo := repository.NewProductRepository(db)

    // Initialize services
    authService := services.NewAuthService(userRepo, cfg.JWTSecret)
    userService := services.NewUserService(userRepo)
    productService := services.NewProductService(productRepo)

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(authService)
    userHandler := handlers.NewUserHandler(userService)
    productHandler := handlers.NewProductHandler(productService)
    fileHandler := handlers.NewFileHandler(cfg.UploadDir, cfg.MaxFileSize)

    // Initialize Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler:          utils.ErrorHandler,
        ReadTimeout:           cfg.Server.ReadTimeout,
        WriteTimeout:          cfg.Server.WriteTimeout,
        IdleTimeout:          cfg.Server.IdleTimeout,
        DisableStartupMessage: cfg.Environment == "production",
    })

    // Global middleware
    app.Use(recover.New())
    app.Use(utils.LoggerMiddleware())
    app.Use(middleware.SecurityMiddleware())
    app.Use(middleware.CORSMiddleware())
    app.Use(middleware.RateLimitMiddleware())

    // Setup routes
    routes.SetupRoutes(app, authHandler, userHandler, productHandler, fileHandler, cfg.JWTSecret)

    // Initialize background workers
    workerPool := workers.NewWorkerPool(5)
    workerPool.Start()
    defer workerPool.Stop()

    // Initialize cron jobs
    scheduler := jobs.NewScheduler(userRepo)
    scheduler.Start()
    defer scheduler.Stop()

    // Start server
    go func() {
        if err := app.Listen(":" + cfg.Server.Port); err != nil {
            log.Fatal("Failed to start server:", err)
        }
    }()

    utils.Logger.Info("Server started on port " + cfg.Server.Port)

    // Graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    <-c
    utils.Logger.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := app.ShutdownWithContext(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    utils.Logger.Info("Server stopped")
}
```

### API Documentation with Swagger

```go
// Add swagger documentation
go get github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/fiber-swagger
go get github.com/swaggo/files
```

```go
// internal/handlers/swagger.go
package handlers

// @title           Fiber API
// @version         1.0
// @description     This is a sample server for Fiber framework.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body services.LoginRequest true "Login credentials"
// @Success      200 {object} services.AuthResponse
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
    // Implementation here
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body services.RegisterRequest true "User registration data"
// @Success      201 {object} services.AuthResponse
// @Failure      400 {object} map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
    // Implementation here
}

// GetMe godoc
// @Summary      Get current user
// @Description  Get current authenticated user information
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} models.User
// @Failure      401 {object} map[string]string
// @Router       /me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
    // Implementation here
}
```

```go
// Add to routes
import (
    _ "my-fiber-app/docs" // Import generated docs
    "github.com/gofiber/swagger"
)

// In SetupRoutes function
app.Get("/swagger/*", swagger.HandlerDefault)
```

### Makefile for Development

```makefile
# Makefile
.PHONY: build run test clean docker-build docker-run migrate-up migrate-down docs

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Generate swagger docs
docs:
	swag init -g cmd/server/main.go

# Database migrations
migrate-up:
	migrate -path migrations -database "postgres://user:password@localhost:5432/myapp?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:password@localhost:5432/myapp?sslmode=disable" down

# Docker commands
docker-build:
	docker build -t my-fiber-app .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

# Development helpers
dev:
	air

lint:
	golangci-lint run

format:
	go fmt ./...

# Install development dependencies
dev-deps:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest
```

### Performance Monitoring

```go
// internal/middleware/metrics.go
package middleware

import (
    "strconv"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests",
        },
        []string{"method", "path"},
    )
)

func PrometheusMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()

        err := c.Next()

        status := strconv.Itoa(c.Response().StatusCode())

        httpRequestsTotal.WithLabelValues(
            c.Method(),
            c.Path(),
            status,
        ).Inc()

        httpRequestDuration.WithLabelValues(
            c.Method(),
            c.Path(),
        ).Observe(time.Since(start).Seconds())

        return err
    }
}
```

---

## Additional Resources

### Useful Go Packages

- **Database**: `gorm.io/gorm`, `github.com/jmoiron/sqlx`
- **Validation**: `github.com/go-playground/validator/v10`
- **Configuration**: `github.com/spf13/viper`, `github.com/joho/godotenv`
- **Testing**: `github.com/stretchr/testify`, `github.com/DATA-DOG/go-sqlmock`
- **Logging**: `github.com/sirupsen/logrus`, `go.uber.org/zap`
- **HTTP Client**: `github.com/go-resty/resty/v2`
- **UUID**: `github.com/google/uuid`
- **Time**: `github.com/jinzhu/now`
- **Crypto**: `golang.org/x/crypto`
- **Background Jobs**: `github.com/hibiken/asynq`

### Learning Resources

1. **Official Go Documentation**: https://golang.org/doc/
2. **Fiber Documentation**: https://gofiber.io/
3. **Go by Example**: https://gobyexample.com/
4. **Effective Go**: https://golang.org/doc/effective_go.html
5. **Go Concurrency Patterns**: https://golang.org/doc/codewalk/sharemem/
6. **Clean Architecture in Go**: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

### Development Tools

- **IDE**: VS Code with Go extension, GoLand
- **Debugging**: Delve debugger
- **API Testing**: Postman, Insomnia
- **Database Tools**: pgAdmin, DBeaver
- **Monitoring**: Prometheus, Grafana
- **Hot Reload**: Air
- **Linting**: golangci-lint

This comprehensive guide covers everything you need to build production-ready backend applications with Go and Fiber. Start with the basics and gradually implement more advanced features as your application grows. Remember to follow Go best practices and maintain clean, testable code throughout your development process.
