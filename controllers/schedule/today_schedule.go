package schedule

import (
	"net/http"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// GetTodaySchedule handles GET /api/schedules/today
func GetTodaySchedule(c *gin.Context) {
	// Get employee ID from context (set during authentication)
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Get today's date (start and end of the day)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	// Define a struct that matches the SQL query column aliases exactly
	type ScheduleData struct {
		ID           uint      `json:"id"`
		DateSchedule string    `json:"date_schedule"`
		ShiftID      uint      `json:"shift_id"`
		ShiftType    string    `json:"shift_type"`
		ShiftStart   string    `json:"shift_start"`
		ShiftEnd     string    `json:"shift_end"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	var scheduleData ScheduleData

	// Query to get employee's schedule for today only
	query := models.DB.Table("schedules s").
		Joins("JOIN shifts sh ON s.shift_id = sh.id").
		Where("s.employee_id = ? AND DATE(s.date_schedule) = ?", 
			employeeID, today.Format("2006-01-02"))

	// Execute the query and retrieve schedule
	result := query.Select(`
		s.id, 
		s.date_schedule, 
		sh.id as shift_id, 
		sh.type as shift_type, 
		sh.start_time as shift_start, 
		sh.end_time as shift_end,
		s.created_at,
		s.updated_at
	`).First(&scheduleData)

	// Check if no schedule found for today
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(http.StatusOK, gin.H{
				"error":     false,
				"message":   "No schedule found for today",
				"schedule":  nil,
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error retrieving today's schedule: " + result.Error.Error(),
		})
		return
	}

	// Convert date format if needed
	dateStr := scheduleData.DateSchedule
	if t, err := time.Parse("2006-01-02", scheduleData.DateSchedule); err == nil {
		dateStr = t.Format("02-01-2006") // Format to DD-MM-YYYY
	}
	
	// Format the response
	schedule := gin.H{
		"id": scheduleData.ID,
		"date_schedule": dateStr,
		"shift": gin.H{
			"id":         scheduleData.ShiftID,
			"start_time": scheduleData.ShiftStart,
			"end_time":   scheduleData.ShiftEnd,
			"type":       scheduleData.ShiftType,
		},
		"created_at": scheduleData.CreatedAt.Format(time.RFC3339),
		"updated_at": scheduleData.UpdatedAt.Format(time.RFC3339),
	}

	// Return the response
	c.JSON(http.StatusOK, gin.H{
		"error":     false,
		"message":   "Today's schedule retrieved successfully",
		"schedule":  schedule,
	})
}