package position

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func DeletePosition(c *gin.Context) {
	var position models.Position
	if err := models.DB.Where("id = ?", c.Param("id")).First(&position).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Record not found!",
		})
		return
	}

	models.DB.Delete(&position)

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Position Deleted Successfully",
	})
}