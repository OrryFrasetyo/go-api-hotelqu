package attendance

import (
	"net/http"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// GetAttendanceLastThreeDays retrieves attendance data for the current day and the previous 2 days
func GetAttendanceLastThreeDays(c *gin.Context) {
	// Get current date
	today := time.Now()
	
	// Calculate dates (today and 2 days before)
	dates := []string{
		today.Format("2006-01-02"),
		today.AddDate(0, 0, -1).Format("2006-01-02"),
		today.AddDate(0, 0, -2).Format("2006-01-02"),
	}

	// Response structure
	type AttendanceResponse struct {
		ID            uint   `json:"id"`
		Employee      struct {
			Name     string `json:"name"`
			Position string `json:"position"`
		} `json:"employee"`
		Schedule      struct {
			ID           uint   `json:"id"`
			DateSchedule string `json:"date_schedule"`
			Status       string `json:"status"`
			Shift        struct {
				ID        uint   `json:"id"`
				Type      string `json:"type"`
				StartTime string `json:"start_end"`
				EndTime   string `json:"end_time"`
			} `json:"shift"`
		} `json:"schedule"`
		Date           string `json:"date"`
		ClockIn        string `json:"clock_in"`
		ClockOut       string `json:"clock_out"`
		Duration       string `json:"duration"`
		ClockInStatus  string `json:"clock_in_status"`
		ClockOutStatus string `json:"clock_out_status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at"`
	}

	var attendances []models.Attendance
	
	// Query attendances for the last three days with all related data
	result := models.DB.
		Preload("Schedule").
		Preload("Schedule.Employee").
		Preload("Schedule.Employee.Position").
		Preload("Schedule.Employee.Position.Department").
		Preload("Schedule.Shift").
		Where("date IN ?", dates).
		Find(&attendances)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve attendance data: " + result.Error.Error(),
		})
		return
	}

	// Process data for response
	var responseAttendances []AttendanceResponse
	for _, att := range attendances {
		var respAtt AttendanceResponse
		respAtt.ID = att.ID
		respAtt.Date = att.Date
		respAtt.ClockIn = att.ClockIn
		respAtt.ClockOut = att.ClockOut
		respAtt.Duration = att.Duration
		respAtt.ClockInStatus = att.ClockInStatus
		respAtt.ClockOutStatus = att.ClockOutStatus
		respAtt.CreatedAt = att.CreatedAt.Format(time.RFC3339)
		respAtt.UpdatedAt = att.UpdatedAt.Format(time.RFC3339)

		// Set schedule data
		respAtt.Schedule.ID = att.Schedule.ID
		respAtt.Schedule.DateSchedule = att.Schedule.DateSchedule
		respAtt.Schedule.Status = att.Schedule.Status
		
		// Set shift data
		respAtt.Schedule.Shift.ID = att.Schedule.Shift.ID
		respAtt.Schedule.Shift.Type = att.Schedule.Shift.Type
		respAtt.Schedule.Shift.StartTime = att.Schedule.Shift.StartTime
		respAtt.Schedule.Shift.EndTime = att.Schedule.Shift.EndTime

		// Set employee data
		respAtt.Employee.Name = att.Schedule.Employee.Name
		if att.Schedule.Employee.Position.PositionName != "" {
			respAtt.Employee.Position = att.Schedule.Employee.Position.PositionName
		}

		responseAttendances = append(responseAttendances, respAtt)
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"error":      false,
		"message":    "Attendance data retrieved successfully",
		"attendances": responseAttendances,
	})
}