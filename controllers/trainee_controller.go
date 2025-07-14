package controllers

import (
	"gin/database"
	"gin/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TraineeController struct{}

func NewTraineeController() *TraineeController {
	return &TraineeController{}
}

func (tc *TraineeController) CreateTrainee(c *gin.Context) {
	var req models.TraineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation based on trainee type
	if req.TraineeType == models.TraineeTypePaidIntern && req.HourlyRate == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hourly rate is required for paid interns"})
		return
	}

	if req.TraineeType == models.TraineeTypeOJT && req.TotalRequiredHours == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total required hours is required for OJT"})
		return
	}

	trainee := models.Trainee{
		UserID:             req.UserID,
		CompanyID:          req.CompanyID,
		EmployeeID:         req.EmployeeID,
		TraineeType:        req.TraineeType,
		HourlyRate:         req.HourlyRate,
		MaxWeeklyHours:     req.MaxWeeklyHours,
		TotalRequiredHours: req.TotalRequiredHours,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		SchoolName:         req.SchoolName,
		Course:             req.Course,
		YearLevel:          req.YearLevel,
	}

	if err := database.DB.Create(&trainee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create trainee"})
		return
	}

	// Load relationships
	database.DB.Preload("User").Preload("Company").First(&trainee, trainee.ID)

	response := models.TraineeResponse{
		ID:                 trainee.ID,
		UserID:             trainee.UserID,
		CompanyID:          trainee.CompanyID,
		EmployeeID:         trainee.EmployeeID,
		TraineeType:        trainee.TraineeType,
		HourlyRate:         trainee.HourlyRate,
		MaxWeeklyHours:     trainee.MaxWeeklyHours,
		TotalRequiredHours: trainee.TotalRequiredHours,
		StartDate:          trainee.StartDate,
		EndDate:            trainee.EndDate,
		Status:             trainee.Status,
		SchoolName:         trainee.SchoolName,
		Course:             trainee.Course,
		YearLevel:          trainee.YearLevel,
		User: models.UserResponse{
			ID:        trainee.User.ID,
			Email:     trainee.User.Email,
			FirstName: trainee.User.FirstName,
			LastName:  trainee.User.LastName,
			Phone:     trainee.User.Phone,
			Role:      trainee.User.Role,
			IsActive:  trainee.User.IsActive,
		},
		Company: models.CompanyResponse{
			ID:            trainee.Company.ID,
			Name:          trainee.Company.Name,
			Address:       trainee.Company.Address,
			ContactPerson: trainee.Company.ContactPerson,
			ContactEmail:  trainee.Company.ContactEmail,
			ContactPhone:  trainee.Company.ContactPhone,
		},
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Trainee created successfully",
		"trainee": response,
	})
}

func (tc *TraineeController) GetTrainees(c *gin.Context) {
	var trainees []models.Trainee
	query := database.DB.Preload("User").Preload("Company")

	// Filter by company if provided
	if companyID := c.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	// Filter by status if provided
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by trainee type if provided
	if traineeType := c.Query("trainee_type"); traineeType != "" {
		query = query.Where("trainee_type = ?", traineeType)
	}

	if err := query.Find(&trainees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch trainees"})
		return
	}

	var responses []models.TraineeResponse
	for _, trainee := range trainees {
		responses = append(responses, models.TraineeResponse{
			ID:                 trainee.ID,
			UserID:             trainee.UserID,
			CompanyID:          trainee.CompanyID,
			EmployeeID:         trainee.EmployeeID,
			TraineeType:        trainee.TraineeType,
			HourlyRate:         trainee.HourlyRate,
			MaxWeeklyHours:     trainee.MaxWeeklyHours,
			TotalRequiredHours: trainee.TotalRequiredHours,
			StartDate:          trainee.StartDate,
			EndDate:            trainee.EndDate,
			Status:             trainee.Status,
			SchoolName:         trainee.SchoolName,
			Course:             trainee.Course,
			YearLevel:          trainee.YearLevel,
			User: models.UserResponse{
				ID:        trainee.User.ID,
				Email:     trainee.User.Email,
				FirstName: trainee.User.FirstName,
				LastName:  trainee.User.LastName,
				Phone:     trainee.User.Phone,
				Role:      trainee.User.Role,
				IsActive:  trainee.User.IsActive,
			},
			Company: models.CompanyResponse{
				ID:            trainee.Company.ID,
				Name:          trainee.Company.Name,
				Address:       trainee.Company.Address,
				ContactPerson: trainee.Company.ContactPerson,
				ContactEmail:  trainee.Company.ContactEmail,
				ContactPhone:  trainee.Company.ContactPhone,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"trainees": responses,
	})
}

func (tc *TraineeController) GetTrainee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	var trainee models.Trainee
	if err := database.DB.Preload("User").Preload("Company").First(&trainee, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trainee not found"})
		return
	}

	response := models.TraineeResponse{
		ID:                 trainee.ID,
		UserID:             trainee.UserID,
		CompanyID:          trainee.CompanyID,
		EmployeeID:         trainee.EmployeeID,
		TraineeType:        trainee.TraineeType,
		HourlyRate:         trainee.HourlyRate,
		MaxWeeklyHours:     trainee.MaxWeeklyHours,
		TotalRequiredHours: trainee.TotalRequiredHours,
		StartDate:          trainee.StartDate,
		EndDate:            trainee.EndDate,
		Status:             trainee.Status,
		SchoolName:         trainee.SchoolName,
		Course:             trainee.Course,
		YearLevel:          trainee.YearLevel,
		User: models.UserResponse{
			ID:        trainee.User.ID,
			Email:     trainee.User.Email,
			FirstName: trainee.User.FirstName,
			LastName:  trainee.User.LastName,
			Phone:     trainee.User.Phone,
			Role:      trainee.User.Role,
			IsActive:  trainee.User.IsActive,
		},
		Company: models.CompanyResponse{
			ID:            trainee.Company.ID,
			Name:          trainee.Company.Name,
			Address:       trainee.Company.Address,
			ContactPerson: trainee.Company.ContactPerson,
			ContactEmail:  trainee.Company.ContactEmail,
			ContactPhone:  trainee.Company.ContactPhone,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"trainee": response,
	})
}

func (tc *TraineeController) UpdateTrainee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	var req models.TraineeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var trainee models.Trainee
	if err := database.DB.First(&trainee, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trainee not found"})
		return
	}

	// Update fields
	trainee.UserID = req.UserID
	trainee.CompanyID = req.CompanyID
	trainee.EmployeeID = req.EmployeeID
	trainee.TraineeType = req.TraineeType
	trainee.HourlyRate = req.HourlyRate
	trainee.MaxWeeklyHours = req.MaxWeeklyHours
	trainee.TotalRequiredHours = req.TotalRequiredHours
	trainee.StartDate = req.StartDate
	trainee.EndDate = req.EndDate
	trainee.SchoolName = req.SchoolName
	trainee.Course = req.Course
	trainee.YearLevel = req.YearLevel

	if err := database.DB.Save(&trainee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update trainee"})
		return
	}

	// Load relationships
	database.DB.Preload("User").Preload("Company").First(&trainee, trainee.ID)

	response := models.TraineeResponse{
		ID:                 trainee.ID,
		UserID:             trainee.UserID,
		CompanyID:          trainee.CompanyID,
		EmployeeID:         trainee.EmployeeID,
		TraineeType:        trainee.TraineeType,
		HourlyRate:         trainee.HourlyRate,
		MaxWeeklyHours:     trainee.MaxWeeklyHours,
		TotalRequiredHours: trainee.TotalRequiredHours,
		StartDate:          trainee.StartDate,
		EndDate:            trainee.EndDate,
		Status:             trainee.Status,
		SchoolName:         trainee.SchoolName,
		Course:             trainee.Course,
		YearLevel:          trainee.YearLevel,
		User: models.UserResponse{
			ID:        trainee.User.ID,
			Email:     trainee.User.Email,
			FirstName: trainee.User.FirstName,
			LastName:  trainee.User.LastName,
			Phone:     trainee.User.Phone,
			Role:      trainee.User.Role,
			IsActive:  trainee.User.IsActive,
		},
		Company: models.CompanyResponse{
			ID:            trainee.Company.ID,
			Name:          trainee.Company.Name,
			Address:       trainee.Company.Address,
			ContactPerson: trainee.Company.ContactPerson,
			ContactEmail:  trainee.Company.ContactEmail,
			ContactPhone:  trainee.Company.ContactPhone,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trainee updated successfully",
		"trainee": response,
	})
}

func (tc *TraineeController) DeleteTrainee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid trainee ID"})
		return
	}

	var trainee models.Trainee
	if err := database.DB.First(&trainee, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trainee not found"})
		return
	}

	if err := database.DB.Delete(&trainee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete trainee"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Trainee deleted successfully",
	})
}