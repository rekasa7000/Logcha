package services

import (
	"gin/database"
	"gin/models"
	"time"
)

type TimeCalculationService struct{}

func NewTimeCalculationService() *TimeCalculationService {
	return &TimeCalculationService{}
}

// CalculateWeeklySummary calculates weekly summary for a trainee
func (tcs *TimeCalculationService) CalculateWeeklySummary(traineeID uint, weekStart time.Time) (*models.WeeklySummary, error) {
	// Get trainee info
	var trainee models.Trainee
	if err := database.DB.First(&trainee, traineeID).Error; err != nil {
		return nil, err
	}

	// Calculate week end date (Sunday)
	weekEnd := weekStart.AddDate(0, 0, 6)

	// Get time records for the week
	var timeRecords []models.TimeRecord
	if err := database.DB.Where("trainee_id = ? AND date BETWEEN ? AND ?", traineeID, weekStart, weekEnd).Find(&timeRecords).Error; err != nil {
		return nil, err
	}

	// Calculate totals
	totalHours := 0.0
	daysPresent := 0

	for _, record := range timeRecords {
		totalHours += record.TotalHours
		if record.Status == models.TimeRecordStatusPresent || 
		   record.Status == models.TimeRecordStatusHalfDayAM || 
		   record.Status == models.TimeRecordStatusHalfDayPM {
			daysPresent++
		}
	}

	// Calculate billable hours based on trainee type
	var billableHours float64
	var grossPay float64

	switch trainee.TraineeType {
	case models.TraineeTypePaidIntern:
		// Cap at max weekly hours
		billableHours = min(totalHours, float64(trainee.MaxWeeklyHours))
		if trainee.HourlyRate != nil {
			grossPay = billableHours * *trainee.HourlyRate
		}
	case models.TraineeTypeUnpaidIntern:
		// Cap at max weekly hours, no pay
		billableHours = min(totalHours, float64(trainee.MaxWeeklyHours))
		grossPay = 0
	case models.TraineeTypeOJT:
		// No weekly limit, all hours count
		billableHours = totalHours
		grossPay = 0
	}

	summary := &models.WeeklySummary{
		TraineeID:        traineeID,
		WeekStartDate:    weekStart,
		WeekEndDate:      weekEnd,
		TotalHoursWorked: totalHours,
		BillableHours:    billableHours,
		GrossPay:         grossPay,
		DaysPresent:      daysPresent,
	}

	return summary, nil
}

// CalculateMonthlyReport calculates monthly report for a trainee
func (tcs *TimeCalculationService) CalculateMonthlyReport(traineeID uint, month, year int) (*models.MonthlyReport, error) {
	// Get trainee info
	var trainee models.Trainee
	if err := database.DB.First(&trainee, traineeID).Error; err != nil {
		return nil, err
	}

	// Calculate month start and end dates
	monthStart := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, -1)

	// Get time records for the month
	var timeRecords []models.TimeRecord
	if err := database.DB.Where("trainee_id = ? AND date BETWEEN ? AND ?", traineeID, monthStart, monthEnd).Find(&timeRecords).Error; err != nil {
		return nil, err
	}

	// Calculate totals
	totalHours := 0.0
	totalBillableHours := 0.0
	totalGrossPay := 0.0
	daysPresent := 0
	daysAbsent := 0

	// Get all weekly summaries for the month
	var weeklySummaries []models.WeeklySummary
	if err := database.DB.Where("trainee_id = ? AND week_start_date BETWEEN ? AND ?", traineeID, monthStart, monthEnd).Find(&weeklySummaries).Error; err != nil {
		return nil, err
	}

	for _, summary := range weeklySummaries {
		totalHours += summary.TotalHoursWorked
		totalBillableHours += summary.BillableHours
		totalGrossPay += summary.GrossPay
		daysPresent += summary.DaysPresent
	}

	// Calculate absent days
	for _, record := range timeRecords {
		if record.Status == models.TimeRecordStatusAbsent {
			daysAbsent++
		}
	}

	report := &models.MonthlyReport{
		TraineeID:          traineeID,
		Month:              month,
		Year:               year,
		TotalHoursWorked:   totalHours,
		TotalBillableHours: totalBillableHours,
		TotalGrossPay:      totalGrossPay,
		DaysPresent:        daysPresent,
		DaysAbsent:         daysAbsent,
		GeneratedAt:        time.Now(),
	}

	return report, nil
}

// GetOJTProgress calculates OJT progress for a trainee
func (tcs *TimeCalculationService) GetOJTProgress(traineeID uint) (map[string]interface{}, error) {
	// Get trainee info
	var trainee models.Trainee
	if err := database.DB.Preload("User").First(&trainee, traineeID).Error; err != nil {
		return nil, err
	}

	if trainee.TraineeType != models.TraineeTypeOJT {
		return nil, nil
	}

	// Get all time records for the trainee
	var timeRecords []models.TimeRecord
	if err := database.DB.Where("trainee_id = ?", traineeID).Find(&timeRecords).Error; err != nil {
		return nil, err
	}

	// Calculate total hours rendered
	totalRendered := 0.0
	for _, record := range timeRecords {
		totalRendered += record.TotalHours
	}

	requiredHours := 0
	if trainee.TotalRequiredHours != nil {
		requiredHours = *trainee.TotalRequiredHours
	}

	remaining := float64(requiredHours) - totalRendered
	if remaining < 0 {
		remaining = 0
	}

	completionPercentage := 0.0
	if requiredHours > 0 {
		completionPercentage = (totalRendered / float64(requiredHours)) * 100
	}

	progress := map[string]interface{}{
		"trainee_id":            traineeID,
		"trainee_name":          trainee.User.FirstName + " " + trainee.User.LastName,
		"total_required_hours":  requiredHours,
		"hours_rendered":        totalRendered,
		"remaining_hours":       remaining,
		"completion_percentage": completionPercentage,
		"is_completed":          remaining == 0,
	}

	return progress, nil
}

// min helper function
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}