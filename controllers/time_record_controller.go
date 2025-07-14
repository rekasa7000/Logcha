package controllers

import (
	"errors"
	"gin/database"
	"gin/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TimeRecordController struct{}

func NewTimeRecordController() *TimeRecordController {
	return &TimeRecordController{}
}

func (trc *TimeRecordController) CreateTimeRecord(c *gin.Context) {
	var req models.TimeRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation
	if err := trc.validateTimeRecord(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if record already exists for this date
	var existingRecord models.TimeRecord
	if err := database.DB.Where("trainee_id = ? AND date = ?", req.TraineeID, req.Date).First(&existingRecord).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Time record already exists for this date"})
		return
	}

	timeRecord := models.TimeRecord{
		TraineeID: req.TraineeID,
		Date:      req.Date,
		AMTimeIn:  req.AMTimeIn,
		AMTimeOut: req.AMTimeOut,
		PMTimeIn:  req.PMTimeIn,
		PMTimeOut: req.PMTimeOut,
		Notes:     req.Notes,
		Status:    req.Status,
	}

	if err := database.DB.Create(&timeRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create time record"})
		return
	}

	response := models.TimeRecordResponse{
		ID:         timeRecord.ID,
		TraineeID:  timeRecord.TraineeID,
		Date:       timeRecord.Date,
		AMTimeIn:   timeRecord.AMTimeIn,
		AMTimeOut:  timeRecord.AMTimeOut,
		AMHours:    timeRecord.AMHours,
		PMTimeIn:   timeRecord.PMTimeIn,
		PMTimeOut:  timeRecord.PMTimeOut,
		PMHours:    timeRecord.PMHours,
		TotalHours: timeRecord.TotalHours,
		Notes:      timeRecord.Notes,
		Status:     timeRecord.Status,
		CreatedAt:  timeRecord.CreatedAt,
		UpdatedAt:  timeRecord.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Time record created successfully",
		"time_record": response,
	})
}

func (trc *TimeRecordController) GetTimeRecords(c *gin.Context) {
	var timeRecords []models.TimeRecord
	query := database.DB.Preload("Trainee").Preload("Trainee.User").Preload("Trainee.Company")

	// Filter by trainee if provided
	if traineeID := c.Query("trainee_id"); traineeID != "" {
		query = query.Where("trainee_id = ?", traineeID)
	}

	// Filter by date range if provided
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("date DESC").Find(&timeRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch time records"})
		return
	}

	var responses []models.TimeRecordResponse
	for _, record := range timeRecords {
		responses = append(responses, models.TimeRecordResponse{
			ID:         record.ID,
			TraineeID:  record.TraineeID,
			Date:       record.Date,
			AMTimeIn:   record.AMTimeIn,
			AMTimeOut:  record.AMTimeOut,
			AMHours:    record.AMHours,
			PMTimeIn:   record.PMTimeIn,
			PMTimeOut:  record.PMTimeOut,
			PMHours:    record.PMHours,
			TotalHours: record.TotalHours,
			Notes:      record.Notes,
			Status:     record.Status,
			CreatedAt:  record.CreatedAt,
			UpdatedAt:  record.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"time_records": responses,
	})
}

func (trc *TimeRecordController) GetTimeRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time record ID"})
		return
	}

	var timeRecord models.TimeRecord
	if err := database.DB.Preload("Trainee").Preload("Trainee.User").Preload("Trainee.Company").First(&timeRecord, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Time record not found"})
		return
	}

	response := models.TimeRecordResponse{
		ID:         timeRecord.ID,
		TraineeID:  timeRecord.TraineeID,
		Date:       timeRecord.Date,
		AMTimeIn:   timeRecord.AMTimeIn,
		AMTimeOut:  timeRecord.AMTimeOut,
		AMHours:    timeRecord.AMHours,
		PMTimeIn:   timeRecord.PMTimeIn,
		PMTimeOut:  timeRecord.PMTimeOut,
		PMHours:    timeRecord.PMHours,
		TotalHours: timeRecord.TotalHours,
		Notes:      timeRecord.Notes,
		Status:     timeRecord.Status,
		CreatedAt:  timeRecord.CreatedAt,
		UpdatedAt:  timeRecord.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"time_record": response,
	})
}

func (trc *TimeRecordController) UpdateTimeRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time record ID"})
		return
	}

	var req models.TimeRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation
	if err := trc.validateTimeRecord(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var timeRecord models.TimeRecord
	if err := database.DB.First(&timeRecord, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Time record not found"})
		return
	}

	// Check if record is older than 7 days (business rule)
	if time.Since(timeRecord.Date) > 7*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot modify time records older than 7 days"})
		return
	}

	// Update fields
	timeRecord.TraineeID = req.TraineeID
	timeRecord.Date = req.Date
	timeRecord.AMTimeIn = req.AMTimeIn
	timeRecord.AMTimeOut = req.AMTimeOut
	timeRecord.PMTimeIn = req.PMTimeIn
	timeRecord.PMTimeOut = req.PMTimeOut
	timeRecord.Notes = req.Notes
	timeRecord.Status = req.Status

	if err := database.DB.Save(&timeRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update time record"})
		return
	}

	response := models.TimeRecordResponse{
		ID:         timeRecord.ID,
		TraineeID:  timeRecord.TraineeID,
		Date:       timeRecord.Date,
		AMTimeIn:   timeRecord.AMTimeIn,
		AMTimeOut:  timeRecord.AMTimeOut,
		AMHours:    timeRecord.AMHours,
		PMTimeIn:   timeRecord.PMTimeIn,
		PMTimeOut:  timeRecord.PMTimeOut,
		PMHours:    timeRecord.PMHours,
		TotalHours: timeRecord.TotalHours,
		Notes:      timeRecord.Notes,
		Status:     timeRecord.Status,
		CreatedAt:  timeRecord.CreatedAt,
		UpdatedAt:  timeRecord.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Time record updated successfully",
		"time_record": response,
	})
}

func (trc *TimeRecordController) DeleteTimeRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time record ID"})
		return
	}

	var timeRecord models.TimeRecord
	if err := database.DB.First(&timeRecord, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Time record not found"})
		return
	}

	// Check if record is older than 7 days (business rule)
	if time.Since(timeRecord.Date) > 7*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete time records older than 7 days"})
		return
	}

	if err := database.DB.Delete(&timeRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete time record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Time record deleted successfully",
	})
}

func (trc *TimeRecordController) validateTimeRecord(req *models.TimeRecordRequest) error {
	// Check if date is in the future
	if req.Date.After(time.Now()) {
		return errors.New("Cannot enter future dates")
	}

	// Validate AM session
	if req.AMTimeIn != nil && req.AMTimeOut != nil {
		if !req.AMTimeOut.After(*req.AMTimeIn) {
			return errors.New("AM time out must be after AM time in")
		}
	}

	// Validate PM session
	if req.PMTimeIn != nil && req.PMTimeOut != nil {
		if !req.PMTimeOut.After(*req.PMTimeIn) {
			return errors.New("PM time out must be after PM time in")
		}
	}

	// Validate PM time in is after AM time out
	if req.AMTimeOut != nil && req.PMTimeIn != nil {
		if !req.PMTimeIn.After(*req.AMTimeOut) {
			return errors.New("PM time in should be after AM time out")
		}
	}

	// At least one session must be present for status = 'present'
	if req.Status == models.TimeRecordStatusPresent {
		if (req.AMTimeIn == nil || req.AMTimeOut == nil) && (req.PMTimeIn == nil || req.PMTimeOut == nil) {
			return errors.New("At least one complete session (AM or PM) is required for present status")
		}
	}

	return nil
}