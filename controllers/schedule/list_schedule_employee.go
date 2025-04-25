package schedule

import (
	"net/http"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// ListSchedules handles GET /api/schedules
func ListSchedules(c *gin.Context) {
	// Get employee ID from context (set during authentication)
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Get current month and next month date range
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := currentMonth.AddDate(0, 2, 0) // Adding 2 months to get the first day of the month after next month
	
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

	var schedulesData []ScheduleData

	// Query to get employee's schedules for current and next month
	query := models.DB.Table("schedules s").
		Joins("JOIN shifts sh ON s.shift_id = sh.id").
		Where("s.employee_id = ? AND DATE(s.date_schedule) >= ? AND DATE(s.date_schedule) < ?", 
			employeeID, currentMonth.Format("2006-01-02"), nextMonth.Format("2006-01-02"))

	// Execute the query and retrieve schedules
	result := query.Select(`
		s.id, 
		s.date_schedule, 
		sh.id as shift_id, 
		sh.type as shift_type, 
		sh.start_time as shift_start, 
		sh.end_time as shift_end,
		s.created_at,
		s.updated_at
	`).Scan(&schedulesData)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error retrieving schedules: " + result.Error.Error(),
		})
		return
	}

	// Convert the raw data to the expected response format
	schedules := make([]gin.H, 0, len(schedulesData))
	for _, s := range schedulesData {
		// Convert date format if needed
		dateStr := s.DateSchedule
		if t, err := time.Parse("2006-01-02", s.DateSchedule); err == nil {
			dateStr = t.Format("02-01-2006") // Format to DD-MM-YYYY
		}
		
		schedules = append(schedules, gin.H{
			"id": s.ID,
			"date_schedule": dateStr,
			"shift": gin.H{
				"id":         s.ShiftID,
				"start_time": s.ShiftStart,
				"end_time":   s.ShiftEnd,
				"type":       s.ShiftType,
			},
			"created_at": s.CreatedAt.Format(time.RFC3339),
			"updated_at": s.UpdatedAt.Format(time.RFC3339),
		})
	}

	// Return the response
	c.JSON(http.StatusOK, gin.H{
		"error":     false,
		"message":   "Schedule data retrieved successfully",
		"schedules": schedules,
	})
}