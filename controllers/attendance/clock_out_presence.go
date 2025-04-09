package attendance

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

type CheckOutRequest struct {
	ClockOut string `json:"clock_out" binding:"required"`
}

// UpdateAttendance handles employee check-out
func UpdateAttendance(c *gin.Context) {
	employeeId, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	var request CheckOutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid request format",
		})
		return
	}

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

	// Find today's attendance record
	var attendance models.Attendance
	attendanceResult := models.DB.Where("schedule_id = ? AND date = ?", schedule.ID, currentDate).First(&attendance)

	if attendanceResult.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "No check-in record found for today. Please check-in first",
		})
		return
	}

	// Check if already clocked out
	if attendance.ClockOut != "" {
		c.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "You have already checked out today",
		})
		return
	}

	// Validate clock-out time
	clockOutStatus := validateClockOut(request.ClockOut, schedule.Shift.EndTime)

	// Calculate duration between clock-in and clock-out
	duration := calculateDuration(attendance.ClockIn, request.ClockOut)

	// Update attendance record
	attendance.ClockOut = request.ClockOut
	attendance.ClockOutStatus = clockOutStatus
	attendance.Duration = duration

	// Use a partial update to avoid overwriting the date field with an incorrect format
	if err := models.DB.Model(&attendance).Updates(map[string]interface{}{
		"clock_out":        request.ClockOut,
		"clock_out_status": clockOutStatus,
		"duration":         duration,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update attendance record: " + err.Error(),
		})
		return
	}

	// Load relations for response
	models.DB.Preload("Schedule").Preload("Schedule.Employee").Preload("Schedule.Employee.Position").Preload("Schedule.Shift").First(&attendance, attendance.ID)

	// Prepare response
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Check-out successful",
		"attendance": gin.H{
			"id": attendance.ID,
			"employee": gin.H{
				"name":     schedule.Employee.Name,
				"position": schedule.Employee.Position.PositionName,
			},
			"schedule": gin.H{
				"id":            schedule.ID,
				"date_schedule": schedule.DateSchedule,
				"status":        schedule.Status,
				"shift": gin.H{
					"id":         schedule.Shift.ID,
					"type":       schedule.Shift.Type,
					"start_time": schedule.Shift.StartTime,
					"end_time":   schedule.Shift.EndTime,
				},
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

// Helper function to validate clock out time and determine status
func validateClockOut(clockOut string, scheduleEnd string) string {
	// Parse clock out and schedule end times (format: HH:MM)
	clockOutParts := strings.Split(clockOut, ":")
	scheduleEndParts := strings.Split(scheduleEnd, ":")

	if len(clockOutParts) < 2 || len(scheduleEndParts) < 2 {
		return "Invalid Format"
	}

	// Convert to integers
	clockOutHour, _ := strconv.Atoi(clockOutParts[0])
	clockOutMinute, _ := strconv.Atoi(clockOutParts[1])
	scheduleEndHour, _ := strconv.Atoi(scheduleEndParts[0])
	scheduleEndMinute, _ := strconv.Atoi(scheduleEndParts[1])

	// Convert times to minutes for easier comparison
	clockOutTotalMinutes := clockOutHour*60 + clockOutMinute
	scheduleEndTotalMinutes := scheduleEndHour*60 + scheduleEndMinute

	// Handle night shift scenario
	isNightShift := isNightShiftSchedule(scheduleEnd)

	if isNightShift {
		// For night shift, if schedule end is earlier in the day (e.g., 06:00),
		// it means it's the next day. We need to adjust our comparison.
		if scheduleEndTotalMinutes < 720 { // Before 12:00 (noon)
			// If clock out is after midnight but before end time, it's on time
			if clockOutTotalMinutes < 720 && clockOutTotalMinutes >= scheduleEndTotalMinutes {
				return "Tepat Waktu"
			} else if clockOutTotalMinutes < scheduleEndTotalMinutes {
				return "Pulang Lebih Awal"
			}
		}
	} else {
		// Normal day shift logic
		if clockOutTotalMinutes < scheduleEndTotalMinutes {
			return "Pulang Lebih Awal"
		}
	}

	// If we reach here, it's on time or late (which is fine for checkout)
	return "Tepat Waktu"
}

// Helper function to check if the schedule is a night shift
func isNightShiftSchedule(endTime string) bool {
	// If end time is in the early morning (e.g., 00:00-07:00),
	// it's likely a night shift
	endParts := strings.Split(endTime, ":")
	endHour, _ := strconv.Atoi(endParts[0])

	return endHour < 7
}

// Helper function to calculate duration between clock-in and clock-out
func calculateDuration(clockIn, clockOut string) string {
	// Parse clock in and clock out times
	inParts := strings.Split(clockIn, ":")
	outParts := strings.Split(clockOut, ":")

	if len(inParts) < 2 || len(outParts) < 2 {
		return "Invalid Format"
	}

	inHour, _ := strconv.Atoi(inParts[0])
	inMinute, _ := strconv.Atoi(inParts[1])
	outHour, _ := strconv.Atoi(outParts[0])
	outMinute, _ := strconv.Atoi(outParts[1])

	// Convert to total minutes
	inTotalMinutes := inHour*60 + inMinute
	outTotalMinutes := outHour*60 + outMinute

	// Handle overnight shifts
	if outTotalMinutes < inTotalMinutes {
		// Add 24 hours (1440 minutes) to out time
		outTotalMinutes += 1440
	}

	// Calculate duration in minutes
	durationMinutes := outTotalMinutes - inTotalMinutes

	// Convert back to hours and minutes
	durationHours := durationMinutes / 60
	remainingMinutes := durationMinutes % 60

	// Format as "X jam Y menit"
	return fmt.Sprintf("%d jam %d menit", durationHours, remainingMinutes)
}