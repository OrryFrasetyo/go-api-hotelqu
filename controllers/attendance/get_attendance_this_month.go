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
				StartTime string `json:"start_time"`
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

	// First get the employee's schedules for this month
	var schedules []models.Schedule
	if err := models.DB.
		Preload("Employee").
		Preload("Employee.Position").
		Preload("Shift").
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
			"error":     false,
			"message":   "No schedules found for this month",
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

	// Create a map of schedules for quick lookup
	scheduleMap := make(map[uint]models.Schedule)
	for _, schedule := range schedules {
		scheduleMap[schedule.ID] = schedule
	}

	// Process data for response
	var responseAttendances []AttendanceResponse
	for _, att := range attendances {
		var respAtt AttendanceResponse
		schedule := scheduleMap[att.ScheduleID]

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
		respAtt.Schedule.ID = schedule.ID
		respAtt.Schedule.DateSchedule = schedule.DateSchedule
		respAtt.Schedule.Status = schedule.Status
		
		// Set shift data
		respAtt.Schedule.Shift.ID = schedule.Shift.ID
		respAtt.Schedule.Shift.Type = schedule.Shift.Type
		respAtt.Schedule.Shift.StartTime = schedule.Shift.StartTime
		respAtt.Schedule.Shift.EndTime = schedule.Shift.EndTime

		// Set employee data
		respAtt.Employee.Name = schedule.Employee.Name
		if schedule.Employee.Position.PositionName != "" {
			respAtt.Employee.Position = schedule.Employee.Position.PositionName
		}

		responseAttendances = append(responseAttendances, respAtt)
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"error":      false,
		"message":    "Attendance data for this month retrieved successfully",
		"attendances": responseAttendances,
	})
}