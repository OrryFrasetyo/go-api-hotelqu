package department

import (
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func FindDepartments(c *gin.Context) {
	var departments []models.Department
	models.DB.Find(&departments)

	c.JSON(200, gin.H{
		"error":   false,
		"message": "List Data Departments",
		"data":    departments,
	})
}