package attendance

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// GetAttendanceByStatus fetches attendance records based on clock_in_status or clock_out_status
func GetAttendanceByStatus(c *gin.Context) {
	// Get employeeId from context (set by JWT middleware)
	employeeId, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Get status params
	clockInStatus := c.Query("clock_in_status")
	clockOutStatus := c.Query("clock_out_status")

	// Validate that only one status parameter is provided
	if (clockInStatus != "" && clockOutStatus != "") || (clockInStatus == "" && clockOutStatus == "") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Please provide either clock_in_status or clock_out_status parameter, not both or none",
		})
		return
	}

	// Initialize database query
	var attendances []models.Attendance
	query := models.DB.
		Preload("Schedule").
		Preload("Schedule.Employee").
		Preload("Schedule.Employee.Position").
		Preload("Schedule.Shift").
		Joins("JOIN schedules ON attendances.schedule_id = schedules.id").
		Where("schedules.employee_id = ?", employeeId)

	// Apply status filter
	if clockInStatus != "" {
		query = query.Where("attendances.clock_in_status = ?", clockInStatus)
	} else if clockOutStatus != "" {
		query = query.Where("attendances.clock_out_status = ?", clockOutStatus)
	}

	// Execute query
	result := query.Find(&attendances)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to fetch attendance records: " + result.Error.Error(),
		})
		return
	}

	// Format response
	formattedAttendances := formatAttendanceResponse(attendances)

	c.JSON(http.StatusOK, gin.H{
		"error":       false,
		"message":     "Attendance data retrieved successfully",
		"attendances": formattedAttendances,
	})
}

// formatAttendanceResponse formats the attendance records for API response
func formatAttendanceResponse(attendances []models.Attendance) []gin.H {
	var formattedAttendances []gin.H

	for _, attendance := range attendances {
		// Format each attendance record according to API spec
		formattedAttendance := gin.H{
			"id": attendance.ID,
			"employee": gin.H{
				"name":     attendance.Schedule.Employee.Name,
				"position": attendance.Schedule.Employee.Position.PositionName,
			},
			"schedule": gin.H{
				"id":            attendance.Schedule.ID,
				"date_schedule": attendance.Schedule.DateSchedule,
				"status":        attendance.Schedule.Status,
				"shift": gin.H{
					"id":         attendance.Schedule.Shift.ID,
					"type":       attendance.Schedule.Shift.Type,
					"start_time": attendance.Schedule.Shift.StartTime,
					"end_time":   attendance.Schedule.Shift.EndTime,
				},
			},
			"date":            attendance.Date,
			"clock_in":        attendance.ClockIn,
			"clock_out":       attendance.ClockOut,
			"duration":        attendance.Duration,
			"clock_in_status": attendance.ClockInStatus,
			"clock_out_status": attendance.ClockOutStatus,
			"created_at":      attendance.CreatedAt,
			"updated_at":      attendance.UpdatedAt,
		}

		formattedAttendances = append(formattedAttendances, formattedAttendance)
	}

	return formattedAttendances
}