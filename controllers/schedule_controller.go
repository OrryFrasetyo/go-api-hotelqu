package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ListDepartmentSchedules handles GET /api/schedules/department
func ListDepartmentSchedules(c *gin.Context) {
	// Extract employee ID from JWT context
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Get current employee with position
	var employee models.Employee
	if err := models.DB.Preload("Position").Preload("Position.Department").First(&employee, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	// Check if employee is manager/supervisor
	var isManager bool
	if result := models.DB.Raw("SELECT position_name LIKE '%manager%' OR position_name LIKE '%supervisor%' FROM positions WHERE id = ?", employee.PositionId).Scan(&isManager); result.Error != nil || !isManager {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You don't have permission to access this resource",
		})
		return
	}

	// Get department ID from employee's position or from query parameter
	departmentID := c.Query("department_id")
	var deptID uint

	if departmentID != "" {
		id, err := strconv.ParseUint(departmentID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid department ID",
			})
			return
		}
		deptID = uint(id)
	} else {
		// Get department ID from employee's position
		deptID = uint(employee.Position.DepartmentId)
	}

	// Get date parameter (optional)
	dateParam := c.Query("date")
	var scheduleDate time.Time
	var err error

	if dateParam != "" {
		// Parse date in DD-MM-YYYY format
		scheduleDate, err = time.Parse("02-01-2006", dateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid date format. Use DD-MM-YYYY",
			})
			return
		}
	} else {
		// Default to today
		scheduleDate = time.Now()
	}

	// Format date as YYYY-MM-DD for database query
	formattedDate := scheduleDate.Format("2006-01-02")
	
	// Get status parameter (optional)
	status := c.Query("status")

	// Get department information
	var department models.Department
	if err := models.DB.First(&department, deptID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Department not found",
		})
		return
	}

	// Count total employees in this department
	var totalEmployees int64
	models.DB.Model(&models.Employee{}).
		Joins("JOIN positions ON employees.position_id = positions.id").
		Where("positions.department_id = ?", deptID).
		Count(&totalEmployees)

	// Define a struct that matches the SQL query column aliases exactly
	type ScheduleData struct {
		ID               uint      `json:"id"`
		EmployeeID       uint      `json:"employee_id"`
		EmployeeName     string    `json:"employee_name"`
		EmployeePosition string    `json:"employee_position"`
		ShiftID          uint      `json:"shift_id"`
		ShiftName        string    `json:"shift_name"`
		ShiftClockIn     string    `json:"shift_clock_in"`
		ShiftClockOut    string    `json:"shift_clock_out"`
		DateSchedule     string    `json:"date_schedule"`
		Status           string    `json:"status"`
		CreatorID        uint      `json:"creator_id"`
		CreatorName      string    `json:"creator_name"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
	}

	var schedulesData []ScheduleData

	// Build query to get schedules from employees in the same department
	query := models.DB.Table("schedules s").
		Joins("JOIN employees e ON s.employee_id = e.id").
		Joins("JOIN positions p ON e.position_id = p.id").
		Joins("JOIN shifts sh ON s.shift_id = sh.id").
		Joins("JOIN employees creator ON s.created_by = creator.id").
		Where("p.department_id = ? AND DATE(s.date_schedule) = ?", deptID, formattedDate)

	// Add status filter if provided
	if status != "" {
		query = query.Where("s.status = ?", status)
	}

	// Execute the query and retrieve schedules
	result := query.Select(`
		s.id, 
		e.id as employee_id, 
		e.name as employee_name, 
		p.position_name as employee_position,
		sh.id as shift_id, 
		sh.type as shift_name, 
		sh.start_time as shift_clock_in, 
		sh.end_time as shift_clock_out,
		s.date_schedule, 
		s.status, 
		creator.id as creator_id, 
		creator.name as creator_name,
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
			dateStr = t.Format("02-01-2006")
		}
		
		schedules = append(schedules, gin.H{
			"id": s.ID,
			"employee": gin.H{
				"id":       s.EmployeeID,
				"name":     s.EmployeeName,
				"position": s.EmployeePosition,
			},
			"shift": gin.H{
				"id":        s.ShiftID,
				"name":      s.ShiftName,
				"clock_in":  s.ShiftClockIn,
				"clock_out": s.ShiftClockOut,
			},
			"date_schedule": dateStr,
			"status":        s.Status,
			"created_by": gin.H{
				"id":   s.CreatorID,
				"name": s.CreatorName,
			},
			"created_at": s.CreatedAt.Format(time.RFC3339),
			"updated_at": s.UpdatedAt.Format(time.RFC3339),
		})
	}

	// Prepare the response
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Schedule data retrieved successfully",
		"meta": gin.H{
			"date": scheduleDate.Format("02-01-2006"),
			"department": gin.H{
				"id":   department.Id,
				"name": department.DepartmentName,
			},
			"total_employees": totalEmployees,
		},
		"schedules": schedules,
	})
}

type CreateScheduleRequest struct {
	EmployeeID   uint   `json:"employee_id" binding:"required"`
	ShiftID      uint   `json:"shift_id" binding:"required"`
	DateSchedule string `json:"date_schedule" binding:"required"` 
	Status       string `json:"status" binding:"required"`
}

func CreateSchedule(c *gin.Context) {
	// Extract employee ID from JWT context
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Get current employee with position
	var creator models.Employee
	if err := models.DB.Preload("Position").Preload("Position.Department").First(&creator, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	// Check if employee is manager/supervisor
	var isManager bool
	if result := models.DB.Raw("SELECT position_name LIKE '%manager%' OR position_name LIKE '%supervisor%' FROM positions WHERE id = ?", creator.PositionId).Scan(&isManager); result.Error != nil || !isManager {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You don't have permission to create schedules",
		})
		return
	}

	// Parse request body
	var request CreateScheduleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		var errors []ErrorMsg
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, ErrorMsg{
				Field:   e.Field(),
				Message: GetErrorMsg(e),
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Validation failed",
			"errors":  errors,
		})
		return
	}

	// Parse date from DD-MM-YYYY format to time.Time
	dateSchedule, err := time.Parse("02-01-2006", request.DateSchedule)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid date format. Use DD-MM-YYYY",
		})
		return
	}

	// Format tanggal ke format yang diterima MySQL (YYYY-MM-DD)
	mysqlFormattedDate := dateSchedule.Format("2006-01-02")

	// Verify employee exists and belongs to the same department as the creator
	var employee models.Employee
	if err := models.DB.Preload("Position").First(&employee, request.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	// Check if employee belongs to the same department as the creator
	if employee.Position.DepartmentId != creator.Position.DepartmentId {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You can only create schedules for employees in your department",
		})
		return
	}

	// Verify shift exists
	var shift models.Shift
	if err := models.DB.First(&shift, request.ShiftID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Shift not found",
		})
		return
	}

	// Check if a schedule already exists for this employee on this date
	var existingSchedule models.Schedule
	result := models.DB.Where("employee_id = ? AND date(date_schedule) = ?", request.EmployeeID, dateSchedule.Format("2006-01-02")).First(&existingSchedule)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "A schedule already exists for this employee on this date",
		})
		return
	}

	// Create new schedule
	schedule := models.Schedule{
		EmployeeID:   request.EmployeeID,
		ShiftID:      request.ShiftID,
		CreatedBy:    uint(creator.Id),
		DateSchedule: mysqlFormattedDate,
		Status:       request.Status,
	}

	// Save to database
	if err := models.DB.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to create schedule: " + err.Error(),
		})
		return
	}

	// Load the schedule with all relationships for response
	var completeSchedule models.Schedule
	if err := models.DB.Preload("Employee").Preload("Employee.Position").Preload("Shift").Preload("Creator").First(&completeSchedule, schedule.ID).Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error":   true,
        "message": "Failed to load schedule details: " + err.Error(),
    })
    return
	}

	// Konversi format tanggal untuk respons
	dateForResponse := completeSchedule.DateSchedule
	// Jika perlu, konversi string dari MySQL format ke format DD-MM-YYYY
	if t, err := time.Parse("2006-01-02", completeSchedule.DateSchedule); err == nil {
			dateForResponse = t.Format("02-01-2006")
	}

	// Format response
	scheduleResponse := gin.H{
    "id": completeSchedule.ID,
    "employee": gin.H{
        "id":       completeSchedule.Employee.Id,
        "name":     completeSchedule.Employee.Name,
        "position": completeSchedule.Employee.Position.PositionName,
    },
    "shift": gin.H{
        "id":        completeSchedule.Shift.ID,
        "name":      completeSchedule.Shift.Type,
        "clock_in":  completeSchedule.Shift.StartTime,
        "clock_out": completeSchedule.Shift.EndTime,
    },
    "date_schedule": dateForResponse,
    "status":        completeSchedule.Status,
    "created_by": gin.H{
        "id":   completeSchedule.Creator.Id,
        "name": completeSchedule.Creator.Name,
    },
    "created_at": completeSchedule.CreatedAt.Format(time.RFC3339),
    "updated_at": completeSchedule.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":    false,
		"message":  "Schedule created successfully",
		"schedule": scheduleResponse,
	})
}

// UpdateSchedule handles PUT /api/schedules/{id}
func UpdateSchedule(c *gin.Context) {
	// Extract employee ID from JWT context
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Get current employee with position
	var employee models.Employee
	if err := models.DB.Preload("Position").Preload("Position.Department").First(&employee, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	// Check if employee is manager/supervisor
	var isManager bool
	if result := models.DB.Raw("SELECT position_name LIKE '%manager%' OR position_name LIKE '%supervisor%' FROM positions WHERE id = ?", employee.PositionId).Scan(&isManager); result.Error != nil || !isManager {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You don't have permission to update schedules",
		})
		return
	}

	// Get schedule ID from URL parameter
	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Schedule ID is required",
		})
		return
	}

	// Validate and convert schedule ID
	id, err := strconv.ParseUint(scheduleID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid schedule ID",
		})
		return
	}

	// Parse request body
	type UpdateScheduleRequest struct {
		ShiftID      uint   `json:"shift_id"`
		DateSchedule string `json:"date_schedule"`
		Status       string `json:"status"`
	}

	var request UpdateScheduleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		var errors []ErrorMsg
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, ErrorMsg{
				Field:   e.Field(),
				Message: GetErrorMsg(e),
			})
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Validation failed",
			"errors":  errors,
		})
		return
	}

	// Find the existing schedule
	var schedule models.Schedule
	if err := models.DB.Preload("Employee").Preload("Employee.Position").First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Schedule not found",
		})
		return
	}

	// Check if the employee belongs to the same department as the manager/supervisor
	if schedule.Employee.Position.DepartmentId != employee.Position.DepartmentId {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You can only update schedules for employees in your department",
		})
		return
	}

	// Update fields
	// Update shift if provided
	if request.ShiftID != 0 {
		// Verify shift exists
		var shift models.Shift
		if err := models.DB.First(&shift, request.ShiftID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "Shift not found",
			})
			return
		}
		schedule.ShiftID = request.ShiftID
	}

	// Update date_schedule if provided
	if request.DateSchedule != "" {
		// Parse date from DD-MM-YYYY format to time.Time
		dateSchedule, err := time.Parse("02-01-2006", request.DateSchedule)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid date format. Use DD-MM-YYYY",
			})
			return
		}

		// Format date to MySQL format (YYYY-MM-DD)
		mysqlFormattedDate := dateSchedule.Format("2006-01-02")
		
		// Check if a schedule already exists for this employee on this date (and it's not this schedule)
		var existingSchedule models.Schedule
		result := models.DB.Where("employee_id = ? AND date(date_schedule) = ? AND id != ?", 
			schedule.EmployeeID, mysqlFormattedDate, schedule.ID).First(&existingSchedule)
		if result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"error":   true,
				"message": "Another schedule already exists for this employee on this date",
			})
			return
		}
		
		schedule.DateSchedule = mysqlFormattedDate
	}

	// Update status if provided
	if request.Status != "" {
		schedule.Status = request.Status
	}

	// Save to database
	if err := models.DB.Save(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update schedule: " + err.Error(),
		})
		return
	}

	// Load the schedule with all relationships for response
	var completeSchedule models.Schedule
	if err := models.DB.Preload("Employee").Preload("Employee.Position").Preload("Shift").Preload("Creator").First(&completeSchedule, schedule.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to load schedule details: " + err.Error(),
		})
		return
	}

	// Convert date format for response (from YYYY-MM-DD to DD-MM-YYYY)
	dateForResponse := completeSchedule.DateSchedule
	if t, err := time.Parse("2006-01-02", completeSchedule.DateSchedule); err == nil {
		dateForResponse = t.Format("02-01-2006")
	}

	// Format response
	scheduleResponse := gin.H{
		"id": completeSchedule.ID,
		"employee": gin.H{
			"id":       completeSchedule.Employee.Id,
			"name":     completeSchedule.Employee.Name,
			"position": completeSchedule.Employee.Position.PositionName,
		},
		"shift": gin.H{
			"id":        completeSchedule.Shift.ID,
			"name":      completeSchedule.Shift.Type,
			"clock_in":  completeSchedule.Shift.StartTime,
			"clock_out": completeSchedule.Shift.EndTime,
		},
		"date_schedule": dateForResponse,
		"status":        completeSchedule.Status,
		"created_by": gin.H{
			"id":   completeSchedule.Creator.Id,
			"name": completeSchedule.Creator.Name,
		},
		"created_at": completeSchedule.CreatedAt.Format(time.RFC3339),
		"updated_at": completeSchedule.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, gin.H{
		"error":    false,
		"message":  "Schedule updated successfully",
		"schedule": scheduleResponse,
	})
}

// DeleteSchedule handles DELETE /api/schedules/{id}
func DeleteSchedule(c *gin.Context) {
	// Extract employee ID from JWT context
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Get current employee with position
	var employee models.Employee
	if err := models.DB.Preload("Position").Preload("Position.Department").First(&employee, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	// Check if employee is manager/supervisor
	var isManager bool
	if result := models.DB.Raw("SELECT position_name LIKE '%manager%' OR position_name LIKE '%supervisor%' FROM positions WHERE id = ?", employee.PositionId).Scan(&isManager); result.Error != nil || !isManager {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You don't have permission to delete schedules",
		})
		return
	}

	// Get schedule ID from URL parameter
	scheduleID := c.Param("id")
	if scheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Schedule ID is required",
		})
		return
	}

	// Validate and convert schedule ID
	id, err := strconv.ParseUint(scheduleID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid schedule ID",
		})
		return
	}

	// Find the existing schedule
	var schedule models.Schedule
	if err := models.DB.Preload("Employee").Preload("Employee.Position").First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Schedule not found",
		})
		return
	}

	// Check if the employee belongs to the same department as the manager/supervisor
	if schedule.Employee.Position.DepartmentId != employee.Position.DepartmentId {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "You can only delete schedules for employees in your department",
		})
		return
	}

	// Delete the schedule
	if err := models.DB.Delete(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to delete schedule: " + err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Schedule deleted successfully",
	})
}