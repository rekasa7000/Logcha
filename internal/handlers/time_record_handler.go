package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/rekasa7000/Logcha/internal/services"
)

type TimeRecordHandler struct {
	timeRecordService *services.TimeRecordService
}

func NewTimeRecordHandler(timeRecordService *services.TimeRecordService) *TimeRecordHandler {
	return &TimeRecordHandler{
		timeRecordService: timeRecordService,
	}
}

func (h *TimeRecordHandler) TimeIn(c fiber.Ctx) error {
	var req services.TimeInRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.timeRecordService.TimeIn(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Time in recorded successfully",
		"data":    response,
	})
}

func (h *TimeRecordHandler) TimeOut(c fiber.Ctx) error {
	var req services.TimeOutRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := h.timeRecordService.TimeOut(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Time out recorded successfully",
		"data":    response,
	})
}

func (h *TimeRecordHandler) GetTimeRecords(c fiber.Ctx) error {
	traineeIDStr := c.Params("traineeId")
	traineeID, err := strconv.ParseUint(traineeIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid trainee ID",
		})
	}

	// Check if date range is provided
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var records []*services.TimeRecordResponse

	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid start date format. Use YYYY-MM-DD",
			})
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid end date format. Use YYYY-MM-DD",
			})
		}

		records, err = h.timeRecordService.GetTimeRecordsByTraineeAndDateRange(c.Context(), uint(traineeID), startDate, endDate)
	} else {
		records, err = h.timeRecordService.GetTimeRecordsByTrainee(c.Context(), uint(traineeID))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": records,
	})
}

func (h *TimeRecordHandler) GetTodaysTimeRecord(c fiber.Ctx) error {
	traineeIDStr := c.Params("traineeId")
	traineeID, err := strconv.ParseUint(traineeIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid trainee ID",
		})
	}

	record, err := h.timeRecordService.GetTodaysTimeRecord(c.Context(), uint(traineeID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No time record found for today",
		})
	}

	return c.JSON(fiber.Map{
		"data": record,
	})
}