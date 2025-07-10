package repository

import (
	"context"
	"time"

	"github.com/rekasa7000/Logcha/internal/models"
	"github.com/rekasa7000/Logcha/pkg/database"
)

type timeRecordRepository struct {
	db *database.DB
}

func NewTimeRecordRepository(db *database.DB) TimeRecordRepository {
	return &timeRecordRepository{db: db}
}

func (r *timeRecordRepository) Create(ctx context.Context, record *models.TimeRecord) error {
	return r.db.WithContext(ctx).Create(record).Error
}

func (r *timeRecordRepository) GetByID(ctx context.Context, id uint) (*models.TimeRecord, error) {
	var record models.TimeRecord
	err := r.db.WithContext(ctx).Preload("Trainee").First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *timeRecordRepository) GetByTraineeID(ctx context.Context, traineeID uint) ([]*models.TimeRecord, error) {
	var records []*models.TimeRecord
	err := r.db.WithContext(ctx).Where("trainee_id = ?", traineeID).Order("date DESC").Find(&records).Error
	return records, err
}

func (r *timeRecordRepository) GetByTraineeIDAndDate(ctx context.Context, traineeID uint, date time.Time) (*models.TimeRecord, error) {
	var record models.TimeRecord
	err := r.db.WithContext(ctx).Where("trainee_id = ? AND date = ?", traineeID, date.Format("2006-01-02")).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *timeRecordRepository) GetByTraineeIDAndDateRange(ctx context.Context, traineeID uint, startDate, endDate time.Time) ([]*models.TimeRecord, error) {
	var records []*models.TimeRecord
	err := r.db.WithContext(ctx).Where("trainee_id = ? AND date BETWEEN ? AND ?", traineeID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Order("date").Find(&records).Error
	return records, err
}

func (r *timeRecordRepository) Update(ctx context.Context, record *models.TimeRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r *timeRecordRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.TimeRecord{}, id).Error
}

func (r *timeRecordRepository) List(ctx context.Context, offset, limit int) ([]*models.TimeRecord, error) {
	var records []*models.TimeRecord
	err := r.db.WithContext(ctx).Preload("Trainee").Offset(offset).Limit(limit).Order("date DESC").Find(&records).Error
	return records, err
}