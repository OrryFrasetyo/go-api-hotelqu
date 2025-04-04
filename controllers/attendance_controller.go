package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// Struktur untuk request check-in
type CheckInRequest struct {
	ClockIn string `json:"clock_in" binding:"required"`
}

// CreateAttendance handles employee check-in
func CreateAttendance(c *gin.Context) {
	// Get employee ID from JWT token
	employeeId, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Parse request body
	var request CheckInRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request format",
		})
		return
	}

	// Get current date in YYYY-MM-DD format
	currentDate := time.Now().Format("2006-01-02")

	// Find employee's schedule for today
	var schedule models.Schedule
	result := models.DB.Where("employee_id = ? AND date_schedule = ?", employeeId, currentDate).
		Preload("Shift").Preload("Employee").Preload("Employee.Position").First(&schedule)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Schedule not found for today",
		})
		return
	}

	// Check if attendance already exists
	var existingAttendance models.Attendance
	checkResult := models.DB.Where("schedule_id = ? AND date = ?", schedule.ID, currentDate).First(&existingAttendance)
	
	if checkResult.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "You have already checked in today",
		})
		return
	}

	// Validate clock-in time
	clockInStatus, isValid := validateClockIn(request.ClockIn, schedule.Shift.StartTime)
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Clock-in is only allowed starting from 1 hour before shift start time",
		})
		return
	}

	// Create new attendance record
	attendance := models.Attendance{
		ScheduleID:     schedule.ID,
		Date:           currentDate,
		ClockIn:        request.ClockIn,
		ClockInStatus:  clockInStatus,
	}

	// Save attendance record
	if err := models.DB.Create(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create attendance record",
		})
		return
	}

	// Load relations for response
	models.DB.Preload("Schedule").Preload("Schedule.Employee").Preload("Schedule.Employee.Position").Preload("Schedule.Shift").First(&attendance, attendance.ID)

	// Prepare response
	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Check-in successful",
		"attendance": gin.H{
			"id": attendance.ID,
			"employee": gin.H{
				"id":       schedule.Employee.Id,
				"name":     schedule.Employee.Name,
				"position": schedule.Employee.Position.PositionName,
			},
			"schedule": gin.H{
				"id":            schedule.ID,
				"date_schedule": schedule.DateSchedule,
				"status":        schedule.Status,
			},
			"date":             attendance.Date,
			"clock_in":         attendance.ClockIn,
			"clock_out":        attendance.ClockOut,
			"duration":         attendance.Duration,
			"clock_in_status":  attendance.ClockInStatus,
			"clock_out_status": attendance.ClockOutStatus,
			"created_at":       attendance.CreatedAt,
			"updated_at":       attendance.UpdatedAt,
		},
	})
}

// Helper function to validate clock in time and determine status
func validateClockIn(clockIn string, scheduleStart string) (string, bool) {
	// Parse clock in and schedule start times (format: HH:MM)
	clockInParts := strings.Split(clockIn, ":")
	scheduleParts := strings.Split(scheduleStart, ":")
	
	if len(clockInParts) < 2 || len(scheduleParts) < 2 {
		return "", false
	}
	
	// Convert to integers
	clockInHour, _ := strconv.Atoi(clockInParts[0])
	clockInMinute, _ := strconv.Atoi(clockInParts[1])
	scheduleHour, _ := strconv.Atoi(scheduleParts[0])
	scheduleMinute, _ := strconv.Atoi(scheduleParts[1])
	
	// Convert times to minutes for easier comparison
	clockInTotalMinutes := clockInHour*60 + clockInMinute
	scheduleTotalMinutes := scheduleHour*60 + scheduleMinute
	
	// Calculate one hour before start time
	oneHourBeforeShift := scheduleTotalMinutes - 60
	
	// If clock-in is before one hour before shift start, it's invalid
	if clockInTotalMinutes < oneHourBeforeShift {
		return "", false
	}
	
	// If clock-in is before or equal to shift start time, it's on time
	if clockInTotalMinutes <= scheduleTotalMinutes {
		return "Tepat Waktu", true
	}
	
	// Otherwise, it's late
	return "Terlambat", true
}