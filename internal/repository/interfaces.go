package repository

import (
	"context"
	"time"

	"github.com/rekasa7000/Logcha/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.User, error)
}

type CompanyRepository interface {
	Create(ctx context.Context, company *models.Company) error
	GetByID(ctx context.Context, id uint) (*models.Company, error)
	Update(ctx context.Context, company *models.Company) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Company, error)
}

type TraineeRepository interface {
	Create(ctx context.Context, trainee *models.Trainee) error
	GetByID(ctx context.Context, id uint) (*models.Trainee, error)
	GetByUserID(ctx context.Context, userID uint) (*models.Trainee, error)
	GetByCompanyID(ctx context.Context, companyID uint) ([]*models.Trainee, error)
	Update(ctx context.Context, trainee *models.Trainee) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.Trainee, error)
}

type TimeRecordRepository interface {
	Create(ctx context.Context, record *models.TimeRecord) error
	GetByID(ctx context.Context, id uint) (*models.TimeRecord, error)
	GetByTraineeID(ctx context.Context, traineeID uint) ([]*models.TimeRecord, error)
	GetByTraineeIDAndDate(ctx context.Context, traineeID uint, date time.Time) (*models.TimeRecord, error)
	GetByTraineeIDAndDateRange(ctx context.Context, traineeID uint, startDate, endDate time.Time) ([]*models.TimeRecord, error)
	Update(ctx context.Context, record *models.TimeRecord) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.TimeRecord, error)
}

type WeeklySummaryRepository interface {
	Create(ctx context.Context, summary *models.WeeklySummary) error
	GetByID(ctx context.Context, id uint) (*models.WeeklySummary, error)
	GetByTraineeID(ctx context.Context, traineeID uint) ([]*models.WeeklySummary, error)
	GetByTraineeIDAndWeek(ctx context.Context, traineeID uint, weekStart time.Time) (*models.WeeklySummary, error)
	Update(ctx context.Context, summary *models.WeeklySummary) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.WeeklySummary, error)
}

type MonthlyReportRepository interface {
	Create(ctx context.Context, report *models.MonthlyReport) error
	GetByID(ctx context.Context, id uint) (*models.MonthlyReport, error)
	GetByTraineeID(ctx context.Context, traineeID uint) ([]*models.MonthlyReport, error)
	GetByTraineeIDAndMonth(ctx context.Context, traineeID uint, month, year int) (*models.MonthlyReport, error)
	Update(ctx context.Context, report *models.MonthlyReport) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]*models.MonthlyReport, error)
}