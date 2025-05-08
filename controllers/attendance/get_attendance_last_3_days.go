package attendance

import (
	"net/http"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// GetAttendanceLastThreeDays retrieves attendance data for the current day and the previous 2 days
func GetAttendanceLastThreeDays(c *gin.Context) {
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "User not authenticated",
		})
		return
	}
	// Get current date
	today := time.Now()

	// Calculate dates (today and 2 days before)
	dates := []string{
		today.Format("2006-01-02"),
		today.AddDate(0, 0, -1).Format("2006-01-02"),
		today.AddDate(0, 0, -2).Format("2006-01-02"),
	}

	// Query attendances for the last three days with only necessary relations
	// And filter by the current employee
	var schedules []models.Schedule
	if err := models.DB.Where("employee_id = ?", employeeID).Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve schedules: " + err.Error(),
		})
		return
	}

	// Get schedule IDs for the employee
	var scheduleIDs []uint
	for _, schedule := range schedules {
		scheduleIDs = append(scheduleIDs, schedule.ID)
	}

	// Query attendances using schedule IDs and dates
	var attendances []models.Attendance
	result := models.DB.Preload("Schedule").
		Preload("Schedule.Employee").
		Preload("Schedule.Employee.Position").
		Preload("Schedule.Shift").
		Where("date IN ? AND schedule_id IN ?", dates, scheduleIDs).
		Find(&attendances)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve attendance data: " + result.Error.Error(),
		})
		return
	}

	// If no attendances found, return empty result
	if len(attendances) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"error":      false,
			"message":    "No attendance data found",
			"attendance": []interface{}{},
		})
		return
	}

	// Process data for response
	var responseAttendances []gin.H
	for _, att := range attendances {
		respAtt := gin.H{
			"id":               att.ID,
			"date":             att.Date,
			"clock_in":         att.ClockIn,
			"clock_out":        att.ClockOut,
			"clock_in_status":  att.ClockInStatus,
			"clock_out_status": att.ClockOutStatus,
			"created_at":       att.CreatedAt.Format(time.RFC3339),
			"updated_at":       att.UpdatedAt.Format(time.RFC3339),
		}

		responseAttendances = append(responseAttendances, respAtt)
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"error":      false,
		"message":    "Attendance data retrieved successfully",
		"attendance": responseAttendances,
	})
}
