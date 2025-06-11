package schedule

import (
	"net/http"
	"strconv"
	"time"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UpdateSchedule handles PUT /api/schedules/{id}
func UpdateSchedule(c *gin.Context) {
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Get current employee with position
	var employee models.Employee
	if err := models.DB.Preload("Position").Preload("Position.Department").First(&employee, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	// Check if employee is manager/supervisor
	var isManager bool
	if result := models.DB.Raw("SELECT position_name LIKE '%manager%' OR position_name LIKE '%supervisor%' OR position_name LIKE '%executive%' FROM positions WHERE id = ?", employee.PositionId).Scan(&isManager); result.Error != nil || !isManager {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You don't have permission to update schedules",
		})
		return
	}

	// Get schedule ID from URL parameter
	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Schedule ID is required",
		})
		return
	}

	// Validate and convert schedule ID
	id, err := strconv.ParseUint(scheduleID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid schedule ID",
		})
		return
	}

	type UpdateScheduleRequest struct {
		ShiftID      uint   `json:"shift_id"`
		DateSchedule string `json:"date_schedule"`
		Status       string `json:"status"`
	}

	var request UpdateScheduleRequest
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

	// Find the existing schedule
	var schedule models.Schedule
	if err := models.DB.Preload("Employee").Preload("Employee.Position").First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Schedule not found",
		})
		return
	}

	// Check if the employee belongs to the same department as the manager/supervisor
	if schedule.Employee.Position.DepartmentId != employee.Position.DepartmentId {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You can only update schedules for employees in your department",
		})
		return
	}

	// Update shift if provided
	if request.ShiftID != 0 {
		var shift models.Shift
		if err := models.DB.First(&shift, request.ShiftID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "Shift not found",
			})
			return
		}
		schedule.ShiftID = request.ShiftID
	}

	// Update date_schedule if provided
	if request.DateSchedule != "" {
		dateSchedule, err := time.Parse("02-01-2006", request.DateSchedule)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid date format. Use DD-MM-YYYY",
			})
			return
		}

		mysqlFormattedDate := dateSchedule.Format("2006-01-02")
		
		// Check if a schedule already exists for this employee on this date (and it's not this schedule)
		var existingSchedule models.Schedule
		result := models.DB.Where("employee_id = ? AND date(date_schedule) = ? AND id != ?", 
			schedule.EmployeeID, mysqlFormattedDate, schedule.ID).First(&existingSchedule)
		if result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error":   true,
				"message": "Another schedule already exists for this employee on this date",
			})
			return
		}
		
		schedule.DateSchedule = mysqlFormattedDate
	}

		schedule.Status = request.Status
	
	if err := models.DB.Save(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update schedule: " + err.Error(),
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

	// Convert date format for response (from YYYY-MM-DD to DD-MM-YYYY)
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

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"message":  "Schedule updated successfully",
		"schedule": scheduleResponse,
	})
}