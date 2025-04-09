package position

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func FindPositionById(c *gin.Context) {
	var position models.Position
	if err := models.DB.Preload("Department").Where("id = ?", c.Param("id")).First(&position).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Record not found!",
		})
		return
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Detail Data Position by ID : " + c.Param("id"),
		"data": map[string]interface{}{
			"id":              position.Id,
			"department_id":   position.DepartmentId,
			"department_name": position.Department.DepartmentName,
			"position_name":   position.PositionName,
			"is_completed":    position.IsCompleted,
		},
	})
}