package repository

import (
	"context"

	"github.com/rekasa7000/Logcha/internal/models"
	"github.com/rekasa7000/Logcha/pkg/database"
)

type traineeRepository struct {
	db *database.DB
}

func NewTraineeRepository(db *database.DB) TraineeRepository {
	return &traineeRepository{db: db}
}

func (r *traineeRepository) Create(ctx context.Context, trainee *models.Trainee) error {
	return r.db.WithContext(ctx).Create(trainee).Error
}

func (r *traineeRepository) GetByID(ctx context.Context, id uint) (*models.Trainee, error) {
	var trainee models.Trainee
	err := r.db.WithContext(ctx).Preload("User").Preload("Company").First(&trainee, id).Error
	if err != nil {
		return nil, err
	}
	return &trainee, nil
}

func (r *traineeRepository) GetByUserID(ctx context.Context, userID uint) (*models.Trainee, error) {
	var trainee models.Trainee
	err := r.db.WithContext(ctx).Preload("User").Preload("Company").Where("user_id = ?", userID).First(&trainee).Error
	if err != nil {
		return nil, err
	}
	return &trainee, nil
}

func (r *traineeRepository) GetByCompanyID(ctx context.Context, companyID uint) ([]*models.Trainee, error) {
	var trainees []*models.Trainee
	err := r.db.WithContext(ctx).Preload("User").Where("company_id = ?", companyID).Find(&trainees).Error
	return trainees, err
}

func (r *traineeRepository) Update(ctx context.Context, trainee *models.Trainee) error {
	return r.db.WithContext(ctx).Save(trainee).Error
}

func (r *traineeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Trainee{}, id).Error
}

func (r *traineeRepository) List(ctx context.Context, offset, limit int) ([]*models.Trainee, error) {
	var trainees []*models.Trainee
	err := r.db.WithContext(ctx).Preload("User").Preload("Company").Offset(offset).Limit(limit).Find(&trainees).Error
	return trainees, err
}