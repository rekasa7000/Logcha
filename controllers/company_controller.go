package controllers

import (
	"gin/database"
	"gin/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CompanyController struct{}

func NewCompanyController() *CompanyController {
	return &CompanyController{}
}

func (cc *CompanyController) CreateCompany(c *gin.Context) {
	var req models.CompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company := models.Company{
		Name:          req.Name,
		Address:       req.Address,
		ContactPerson: req.ContactPerson,
		ContactEmail:  req.ContactEmail,
		ContactPhone:  req.ContactPhone,
	}

	if err := database.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	response := models.CompanyResponse{
		ID:            company.ID,
		Name:          company.Name,
		Address:       company.Address,
		ContactPerson: company.ContactPerson,
		ContactEmail:  company.ContactEmail,
		ContactPhone:  company.ContactPhone,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Company created successfully",
		"company": response,
	})
}

func (cc *CompanyController) GetCompanies(c *gin.Context) {
	var companies []models.Company
	if err := database.DB.Find(&companies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch companies"})
		return
	}

	var responses []models.CompanyResponse
	for _, company := range companies {
		responses = append(responses, models.CompanyResponse{
			ID:            company.ID,
			Name:          company.Name,
			Address:       company.Address,
			ContactPerson: company.ContactPerson,
			ContactEmail:  company.ContactEmail,
			ContactPhone:  company.ContactPhone,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"companies": responses,
	})
}

func (cc *CompanyController) GetCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var company models.Company
	if err := database.DB.First(&company, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	response := models.CompanyResponse{
		ID:            company.ID,
		Name:          company.Name,
		Address:       company.Address,
		ContactPerson: company.ContactPerson,
		ContactEmail:  company.ContactEmail,
		ContactPhone:  company.ContactPhone,
	}

	c.JSON(http.StatusOK, gin.H{
		"company": response,
	})
}

func (cc *CompanyController) UpdateCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var req models.CompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var company models.Company
	if err := database.DB.First(&company, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	company.Name = req.Name
	company.Address = req.Address
	company.ContactPerson = req.ContactPerson
	company.ContactEmail = req.ContactEmail
	company.ContactPhone = req.ContactPhone

	if err := database.DB.Save(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	response := models.CompanyResponse{
		ID:            company.ID,
		Name:          company.Name,
		Address:       company.Address,
		ContactPerson: company.ContactPerson,
		ContactEmail:  company.ContactEmail,
		ContactPhone:  company.ContactPhone,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company updated successfully",
		"company": response,
	})
}

func (cc *CompanyController) DeleteCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var company models.Company
	if err := database.DB.First(&company, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	if err := database.DB.Delete(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company deleted successfully",
	})
}