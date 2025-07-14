package controllers

import (
	"gin/database"
	"gin/models"
	"gin/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportsController struct {
	timeCalcService *services.TimeCalculationService
}

func NewReportsController() *ReportsController {
	return &ReportsController{
		timeCalcService: services.NewTimeCalculationService(),
	}
}

// GetWeeklySummary generates and returns weekly summary
func (rc *ReportsController) GetWeeklySummary(c *gin.Context) {
	traineeIDStr := c.Query("trainee_id")
	weekStartStr := c.Query("week_start")

	if traineeIDStr == "" || weekStartStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trainee_id and week_start are required"})
		return
	}

	traineeID, err := strconv.ParseUint(traineeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	weekStart, err := time.Parse("2006-01-02", weekStartStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week_start format. Use YYYY-MM-DD"})
		return
	}

	summary, err := rc.timeCalcService.CalculateWeeklySummary(uint(traineeID), weekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate weekly summary"})
		return
	}

	// Save or update the summary
	var existingSummary models.WeeklySummary
	if err := database.DB.Where("trainee_id = ? AND week_start_date = ?", traineeID, weekStart).First(&existingSummary).Error; err == nil {
		// Update existing
		existingSummary.TotalHoursWorked = summary.TotalHoursWorked
		existingSummary.BillableHours = summary.BillableHours
		existingSummary.GrossPay = summary.GrossPay
		existingSummary.DaysPresent = summary.DaysPresent
		database.DB.Save(&existingSummary)
		summary = &existingSummary
	} else {
		// Create new
		database.DB.Create(summary)
	}

	response := models.WeeklySummaryResponse{
		ID:               summary.ID,
		TraineeID:        summary.TraineeID,
		WeekStartDate:    summary.WeekStartDate,
		WeekEndDate:      summary.WeekEndDate,
		TotalHoursWorked: summary.TotalHoursWorked,
		BillableHours:    summary.BillableHours,
		GrossPay:         summary.GrossPay,
		DaysPresent:      summary.DaysPresent,
		CreatedAt:        summary.CreatedAt,
		UpdatedAt:        summary.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"weekly_summary": response,
	})
}

// GetMonthlyReport generates and returns monthly report
func (rc *ReportsController) GetMonthlyReport(c *gin.Context) {
	traineeIDStr := c.Query("trainee_id")
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if traineeIDStr == "" || monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trainee_id, month, and year are required"})
		return
	}

	traineeID, err := strconv.ParseUint(traineeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month. Must be 1-12"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2020 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}

	report, err := rc.timeCalcService.CalculateMonthlyReport(uint(traineeID), month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate monthly report"})
		return
	}

	// Save or update the report
	var existingReport models.MonthlyReport
	if err := database.DB.Where("trainee_id = ? AND month = ? AND year = ?", traineeID, month, year).First(&existingReport).Error; err == nil {
		// Update existing
		existingReport.TotalHoursWorked = report.TotalHoursWorked
		existingReport.TotalBillableHours = report.TotalBillableHours
		existingReport.TotalGrossPay = report.TotalGrossPay
		existingReport.DaysPresent = report.DaysPresent
		existingReport.DaysAbsent = report.DaysAbsent
		existingReport.GeneratedAt = time.Now()
		database.DB.Save(&existingReport)
		report = &existingReport
	} else {
		// Create new
		database.DB.Create(report)
	}

	response := models.MonthlyReportResponse{
		ID:                 report.ID,
		TraineeID:          report.TraineeID,
		Month:              report.Month,
		Year:               report.Year,
		TotalHoursWorked:   report.TotalHoursWorked,
		TotalBillableHours: report.TotalBillableHours,
		TotalGrossPay:      report.TotalGrossPay,
		DaysPresent:        report.DaysPresent,
		DaysAbsent:         report.DaysAbsent,
		GeneratedAt:        report.GeneratedAt,
		CreatedAt:          report.CreatedAt,
		UpdatedAt:          report.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"monthly_report": response,
	})
}

// GetOJTProgress returns OJT progress for a trainee
func (rc *ReportsController) GetOJTProgress(c *gin.Context) {
	traineeIDStr := c.Param("trainee_id")
	if traineeIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trainee_id is required"})
		return
	}

	traineeID, err := strconv.ParseUint(traineeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	progress, err := rc.timeCalcService.GetOJTProgress(uint(traineeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get OJT progress"})
		return
	}

	if progress == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Trainee is not an OJT student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ojt_progress": progress,
	})
}

// GetAllOJTProgress returns OJT progress for all OJT trainees
func (rc *ReportsController) GetAllOJTProgress(c *gin.Context) {
	// Get all OJT trainees
	var trainees []models.Trainee
	if err := database.DB.Where("trainee_type = ? AND status = ?", models.TraineeTypeOJT, models.TraineeStatusActive).Find(&trainees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch OJT trainees"})
		return
	}

	var progressList []map[string]interface{}
	for _, trainee := range trainees {
		progress, err := rc.timeCalcService.GetOJTProgress(trainee.ID)
		if err != nil {
			continue
		}
		if progress != nil {
			progressList = append(progressList, progress)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"ojt_progress_list": progressList,
	})
}

// GetDTR returns Daily Time Record for a trainee
func (rc *ReportsController) GetDTR(c *gin.Context) {
	traineeIDStr := c.Query("trainee_id")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if traineeIDStr == "" || startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trainee_id, start_date, and end_date are required"})
		return
	}

	traineeID, err := strconv.ParseUint(traineeIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	// Get time records for the period
	var timeRecords []models.TimeRecord
	if err := database.DB.Where("trainee_id = ? AND date BETWEEN ? AND ?", traineeID, startDate, endDate).Order("date ASC").Find(&timeRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch time records"})
		return
	}

	var dtrRecords []map[string]interface{}
	for _, record := range timeRecords {
		dtrRecord := map[string]interface{}{
			"date":        record.Date.Format("2006-01-02"),
			"am_time_in":  record.AMTimeIn,
			"am_time_out": record.AMTimeOut,
			"pm_time_in":  record.PMTimeIn,
			"pm_time_out": record.PMTimeOut,
			"total_hours": record.TotalHours,
			"status":      record.Status,
			"notes":       record.Notes,
		}
		dtrRecords = append(dtrRecords, dtrRecord)
	}

	c.JSON(http.StatusOK, gin.H{
		"dtr_records": dtrRecords,
		"period": map[string]string{
			"start_date": startDate.Format("2006-01-02"),
			"end_date":   endDate.Format("2006-01-02"),
		},
	})
}