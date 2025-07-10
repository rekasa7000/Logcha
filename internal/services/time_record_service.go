package services

import (
	"context"
	"errors"
	"time"

	"github.com/rekasa7000/Logcha/internal/models"
	"github.com/rekasa7000/Logcha/internal/repository"
	"github.com/rekasa7000/Logcha/internal/utils"
)

type TimeRecordService struct {
	timeRecordRepo repository.TimeRecordRepository
	traineeRepo    repository.TraineeRepository
}

func NewTimeRecordService(timeRecordRepo repository.TimeRecordRepository, traineeRepo repository.TraineeRepository) *TimeRecordService {
	return &TimeRecordService{
		timeRecordRepo: timeRecordRepo,
		traineeRepo:    traineeRepo,
	}
}

type TimeInRequest struct {
	TraineeID uint   `json:"trainee_id" validate:"required"`
	Session   string `json:"session" validate:"required,oneof=am pm"`
}

type TimeOutRequest struct {
	TraineeID uint   `json:"trainee_id" validate:"required"`
	Session   string `json:"session" validate:"required,oneof=am pm"`
}

type TimeRecordResponse struct {
	ID         uint      `json:"id"`
	TraineeID  uint      `json:"trainee_id"`
	Date       time.Time `json:"date"`
	AMTimeIn   *time.Time `json:"am_time_in"`
	AMTimeOut  *time.Time `json:"am_time_out"`
	PMTimeIn   *time.Time `json:"pm_time_in"`
	PMTimeOut  *time.Time `json:"pm_time_out"`
	AMHours    float64   `json:"am_hours"`
	PMHours    float64   `json:"pm_hours"`
	TotalHours float64   `json:"total_hours"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
}

func (s *TimeRecordService) TimeIn(ctx context.Context, req *TimeInRequest) (*TimeRecordResponse, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, errors.New("validation failed")
	}

	// Check if trainee exists
	trainee, err := s.traineeRepo.GetByID(ctx, req.TraineeID)
	if err != nil {
		return nil, errors.New("trainee not found")
	}

	// Check if trainee is active
	if trainee.Status != models.TraineeStatusActive {
		return nil, errors.New("trainee is not active")
	}

	today := time.Now()
	dateOnly := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	// Get or create today's time record
	record, err := s.timeRecordRepo.GetByTraineeIDAndDate(ctx, req.TraineeID, dateOnly)
	if err != nil {
		// Create new record if doesn't exist
		record = &models.TimeRecord{
			TraineeID: req.TraineeID,
			Date:      dateOnly,
			Status:    models.TimeRecordStatusPresent,
		}
	}

	// Set time in based on session
	now := time.Now()
	timeOnly := time.Date(0, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	if req.Session == "am" {
		if record.AMTimeIn != nil {
			return nil, errors.New("already timed in for AM session")
		}
		record.AMTimeIn = &timeOnly
	} else {
		if record.PMTimeIn != nil {
			return nil, errors.New("already timed in for PM session")
		}
		record.PMTimeIn = &timeOnly
	}

	// Save or update record
	if record.ID == 0 {
		err = s.timeRecordRepo.Create(ctx, record)
	} else {
		err = s.timeRecordRepo.Update(ctx, record)
	}

	if err != nil {
		return nil, err
	}

	return s.toTimeRecordResponse(record), nil
}

func (s *TimeRecordService) TimeOut(ctx context.Context, req *TimeOutRequest) (*TimeRecordResponse, error) {
	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		return nil, errors.New("validation failed")
	}

	// Check if trainee exists
	trainee, err := s.traineeRepo.GetByID(ctx, req.TraineeID)
	if err != nil {
		return nil, errors.New("trainee not found")
	}

	// Check if trainee is active
	if trainee.Status != models.TraineeStatusActive {
		return nil, errors.New("trainee is not active")
	}

	today := time.Now()
	dateOnly := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	// Get today's time record
	record, err := s.timeRecordRepo.GetByTraineeIDAndDate(ctx, req.TraineeID, dateOnly)
	if err != nil {
		return nil, errors.New("no time in record found for today")
	}

	// Set time out based on session
	now := time.Now()
	timeOnly := time.Date(0, 1, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	if req.Session == "am" {
		if record.AMTimeIn == nil {
			return nil, errors.New("no AM time in record found")
		}
		if record.AMTimeOut != nil {
			return nil, errors.New("already timed out for AM session")
		}
		record.AMTimeOut = &timeOnly
	} else {
		if record.PMTimeIn == nil {
			return nil, errors.New("no PM time in record found")
		}
		if record.PMTimeOut != nil {
			return nil, errors.New("already timed out for PM session")
		}
		record.PMTimeOut = &timeOnly
	}

	// Update record
	err = s.timeRecordRepo.Update(ctx, record)
	if err != nil {
		return nil, err
	}

	return s.toTimeRecordResponse(record), nil
}

func (s *TimeRecordService) GetTimeRecordsByTrainee(ctx context.Context, traineeID uint) ([]*TimeRecordResponse, error) {
	records, err := s.timeRecordRepo.GetByTraineeID(ctx, traineeID)
	if err != nil {
		return nil, err
	}

	var responses []*TimeRecordResponse
	for _, record := range records {
		responses = append(responses, s.toTimeRecordResponse(record))
	}

	return responses, nil
}

func (s *TimeRecordService) GetTimeRecordsByTraineeAndDateRange(ctx context.Context, traineeID uint, startDate, endDate time.Time) ([]*TimeRecordResponse, error) {
	records, err := s.timeRecordRepo.GetByTraineeIDAndDateRange(ctx, traineeID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var responses []*TimeRecordResponse
	for _, record := range records {
		responses = append(responses, s.toTimeRecordResponse(record))
	}

	return responses, nil
}

func (s *TimeRecordService) GetTodaysTimeRecord(ctx context.Context, traineeID uint) (*TimeRecordResponse, error) {
	today := time.Now()
	dateOnly := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	record, err := s.timeRecordRepo.GetByTraineeIDAndDate(ctx, traineeID, dateOnly)
	if err != nil {
		return nil, err
	}

	return s.toTimeRecordResponse(record), nil
}

func (s *TimeRecordService) toTimeRecordResponse(record *models.TimeRecord) *TimeRecordResponse {
	return &TimeRecordResponse{
		ID:         record.ID,
		TraineeID:  record.TraineeID,
		Date:       record.Date,
		AMTimeIn:   record.AMTimeIn,
		AMTimeOut:  record.AMTimeOut,
		PMTimeIn:   record.PMTimeIn,
		PMTimeOut:  record.PMTimeOut,
		AMHours:    record.AMHours,
		PMHours:    record.PMHours,
		TotalHours: record.TotalHours,
		Status:     string(record.Status),
		Notes:      record.Notes,
	}
}