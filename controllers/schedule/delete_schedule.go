package schedule

import (
	"net/http"
	"strconv"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// DeleteSchedule handles DELETE /api/schedules/{id}
func DeleteSchedule(c *gin.Context) {
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
	if result := models.DB.Raw("SELECT position_name LIKE '%manager%' OR position_name LIKE '%supervisor%' OR position_name LIKE '%executive%'  FROM positions WHERE id = ?", employee.PositionId).Scan(&isManager); result.Error != nil || !isManager {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You don't have permission to delete schedules",
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
			"message": "You can only delete schedules for employees in your department",
		})
		return
	}

	// Delete the schedule
	if err := models.DB.Delete(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete schedule: " + err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Schedule deleted successfully",
	})
}