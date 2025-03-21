package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
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

	// Execute the query and retrieve schedules
	type ScheduleResponse struct {
		ID           uint      `json:"id"`
		Employee     struct {
			ID       uint   `json:"id"`
			Name     string `json:"name"`
			Position string `json:"position"`
		} `json:"employee"`
		Shift struct {
			ID       uint   `json:"id"`
			Name     string `json:"name"`
			ClockIn  string `json:"clock_in"`
			ClockOut string `json:"clock_out"`
		} `json:"shift"`
		DateSchedule string `json:"date_schedule"`
		Status       string `json:"status"`
		CreatedBy    struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"created_by"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	var schedules []ScheduleResponse
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
	`).
		Scan(&schedules)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error retrieving schedules: " + result.Error.Error(),
		})
		return
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