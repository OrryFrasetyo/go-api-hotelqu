package attendance

import (
	"net/http"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// GetAttendanceToday retrieves today's attendance data for the currently logged in employee
func GetAttendanceToday(c *gin.Context) {
	// Get employee ID from JWT token (set in middleware)
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Employee ID not found in token",
		})
		return
	}

	// Get current date
	today := time.Now().Format("2006-01-02")

	// Define response structure
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

	// Get the employee's schedule for today
	var schedule models.Schedule
	if err := models.DB.
		Preload("Employee").
		Preload("Employee.Position").
		Preload("Shift").
		Where("employee_id = ? AND date_schedule = ?", employeeID, today).
		First(&schedule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "No schedule found for today",
		})
		return
	}

	// Get attendance record for today's schedule
	var attendance models.Attendance
	result := models.DB.
		Where("schedule_id = ? AND date = ?", schedule.ID, today).
		First(&attendance)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "No attendance record found for today",
		})
		return
	}

	// Prepare response
	var response AttendanceResponse
	response.ID = attendance.ID
	response.Date = attendance.Date
	response.ClockIn = attendance.ClockIn
	response.ClockOut = attendance.ClockOut
	response.Duration = attendance.Duration
	response.ClockInStatus = attendance.ClockInStatus
	response.ClockOutStatus = attendance.ClockOutStatus
	response.CreatedAt = attendance.CreatedAt.Format(time.RFC3339)
	response.UpdatedAt = attendance.UpdatedAt.Format(time.RFC3339)

	// Set schedule data
	response.Schedule.ID = schedule.ID
	response.Schedule.DateSchedule = schedule.DateSchedule
	response.Schedule.Status = schedule.Status
	
	// Set shift data
	response.Schedule.Shift.ID = schedule.Shift.ID
	response.Schedule.Shift.Type = schedule.Shift.Type
	response.Schedule.Shift.StartTime = schedule.Shift.StartTime
	response.Schedule.Shift.EndTime = schedule.Shift.EndTime

	// Set employee data
	response.Employee.Name = schedule.Employee.Name
	if schedule.Employee.Position.PositionName != "" {
		response.Employee.Position = schedule.Employee.Position.PositionName
	}

	// Return response
	c.JSON(http.StatusOK, gin.H{
		"error":          false,
		"message":        "Today's attendance data retrieved successfully",
		"attendance_now": response,
	})
}