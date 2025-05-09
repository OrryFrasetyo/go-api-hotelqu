package attendance

import (
	"net/http"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// GetAttendanceThisMonth retrieves attendance data for the current month for logged in employee
func GetAttendanceThisMonth(c *gin.Context) {
	// Get employee ID from JWT token (set in middleware)
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Employee ID not found in token",
		})
		return
	}

	// Get current date info
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()

	// Get first day of the month
	firstDay := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())

	// Get last day of the month
	lastDay := firstDay.AddDate(0, 1, -1)

	// Format dates for SQL query
	startDate := firstDay.Format("2006-01-02")
	endDate := lastDay.Format("2006-01-02")

	// Response structure
	type AttendanceResponse struct {
		ID             uint   `json:"id"`
		Date           string `json:"date"`
		ClockIn        string `json:"clock_in"`
		ClockOut       string `json:"clock_out"`
		ClockInStatus  string `json:"clock_in_status"`
		ClockOutStatus string `json:"clock_out_status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at"`
	}

	// First get the employee's schedules for this month
	var schedules []models.Schedule
	if err := models.DB.
		Where("employee_id = ? AND date_schedule BETWEEN ? AND ?", employeeID, startDate, endDate).
		Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve schedules: " + err.Error(),
		})
		return
	}

	// Get all schedule IDs
	var scheduleIDs []uint
	for _, schedule := range schedules {
		scheduleIDs = append(scheduleIDs, schedule.ID)
	}

	// If no schedules found
	if len(scheduleIDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error":       false,
			"message":     "No schedules found for this month",
			"attendances": []interface{}{},
		})
		return
	}

	// Get attendances for these schedules
	var attendances []models.Attendance
	if err := models.DB.
		Where("schedule_id IN ? AND date BETWEEN ? AND ?", scheduleIDs, startDate, endDate).
		Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve attendance data: " + err.Error(),
		})
		return
	}

	// Process data for response - simplified according to API spec
	var responseAttendances []AttendanceResponse
	for _, att := range attendances {
		responseAttendances = append(responseAttendances, AttendanceResponse{
			ID:             att.ID,
			Date:           att.Date,
			ClockIn:        att.ClockIn,
			ClockOut:       att.ClockOut,
			ClockInStatus:  att.ClockInStatus,
			ClockOutStatus: att.ClockOutStatus,
			CreatedAt:      att.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      att.UpdatedAt.Format(time.RFC3339),
		})
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"error":       false,
		"message":     "Attendance data for this month retrieved successfully",
		"attendances": responseAttendances,
	})
}
