package schedule

import (
	"net/http"
	"time"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CreateScheduleRequest struct {
	EmployeeID   uint   `json:"employee_id" binding:"required"`
	ShiftID      uint   `json:"shift_id" binding:"required"`
	DateSchedule string `json:"date_schedule" binding:"required"`
	Status       string `json:"status" binding:"required"`
}

func CreateSchedule(c *gin.Context) {
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Get current employee with position
	var creator models.Employee
	if err := models.DB.Preload("Position").Preload("Position.Department").First(&creator, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	// Check if employee is manager/supervisor
	var isManager bool
	if result := models.DB.Raw("SELECT position_name LIKE '%manager%' OR position_name LIKE '%supervisor%' OR position_name LIKE '%executive%' FROM positions WHERE id = ?", creator.PositionId).Scan(&isManager); result.Error != nil || !isManager {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You don't have permission to create schedules",
		})
		return
	}

	// Parse request body
	var request CreateScheduleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		var errors []errormessage.ErrorMsg
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, errormessage.ErrorMsg{
				Field:   e.Field(),
				Message: errormessage.GetErrorMsg(e),
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Validation failed",
			"errors":  errors,
		})
		return
	}

	// Parse date from DD-MM-YYYY format to time.Time
	dateSchedule, err := time.Parse("02-01-2006", request.DateSchedule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid date format. Use DD-MM-YYYY",
		})
		return
	}

	mysqlFormattedDate := dateSchedule.Format("2006-01-02")

	// Verify employee exists and belongs to the same department as the creator
	var employee models.Employee
	if err := models.DB.Preload("Position").First(&employee, request.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	// Check if employee belongs to the same department as the creator
	if employee.Position.DepartmentId != creator.Position.DepartmentId {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You can only create schedules for employees in your department",
		})
		return
	}

	// Verify shift exists
	var shift models.Shift
	if err := models.DB.First(&shift, request.ShiftID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Shift not found",
		})
		return
	}

	// Check if a schedule already exists for this employee on this date
	var existingSchedule models.Schedule
	result := models.DB.Where("employee_id = ? AND date(date_schedule) = ?", request.EmployeeID, dateSchedule.Format("2006-01-02")).First(&existingSchedule)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "A schedule already exists for this employee on this date",
		})
		return
	}

	schedule := models.Schedule{
		EmployeeID:   request.EmployeeID,
		ShiftID:      request.ShiftID,
		CreatedBy:    uint(creator.Id),
		DateSchedule: mysqlFormattedDate,
		Status:       request.Status,
	}

	if err := models.DB.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create schedule: " + err.Error(),
		})
		return
	}

	// Load the schedule with all relationships for response
	var completeSchedule models.Schedule
	if err := models.DB.Preload("Employee").Preload("Employee.Position").Preload("Shift").Preload("Creator").First(&completeSchedule, schedule.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to load schedule details: " + err.Error(),
		})
		return
	}

	dateForResponse := completeSchedule.DateSchedule
	if t, err := time.Parse("2006-01-02", completeSchedule.DateSchedule); err == nil {
		dateForResponse = t.Format("02-01-2006")
	}

	// Format response
	scheduleResponse := gin.H{
		"id": completeSchedule.ID,
		"employee": gin.H{
			"id":       completeSchedule.Employee.Id,
			"name":     completeSchedule.Employee.Name,
			"position": completeSchedule.Employee.Position.PositionName,
		},
		"shift": gin.H{
			"id":        completeSchedule.Shift.ID,
			"name":      completeSchedule.Shift.Type,
			"clock_in":  completeSchedule.Shift.StartTime,
			"clock_out": completeSchedule.Shift.EndTime,
		},
		"date_schedule": dateForResponse,
		"status":        completeSchedule.Status,
		"created_by": gin.H{
			"id":   completeSchedule.Creator.Id,
			"name": completeSchedule.Creator.Name,
		},
		"created_at": completeSchedule.CreatedAt.Format(time.RFC3339),
		"updated_at": completeSchedule.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":    false,
		"message":  "Schedule created successfully",
		"schedule": scheduleResponse,
	})
}