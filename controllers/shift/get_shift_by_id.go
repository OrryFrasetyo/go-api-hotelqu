package shift

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func FindShiftById(c *gin.Context) {
	var shift models.Shift
	id := c.Param("id")

	if err := models.DB.First(&shift, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "true",
			"message": "Shift not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   "false",
		"message": "Shift found",
		"data":    shift,
	})
}