package employee

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	employeeId, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	var employee models.Employee
	if err := models.DB.Preload("Position.Department").First(&employee, employeeId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	var photoURL *string
	if employee.Photo != nil {
		photoURL = employee.Photo
	}

	var departmentName string
	if employee.Position.Department.DepartmentName != "" {
		departmentName = employee.Position.Department.DepartmentName
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Profile retrieved successfully",
		"profile": gin.H{
			"id":         employee.Id,
			"name":       employee.Name,
			"email":      employee.Email,
			"phone":      employee.Phone,
			"position":   employee.Position.PositionName,
			"department": departmentName,
			"photo":      photoURL,
		},
	})
}