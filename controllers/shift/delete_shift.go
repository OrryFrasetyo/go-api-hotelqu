package shift

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func DeleteShift(c *gin.Context) {
	var shift models.Shift
	id := c.Param("id")

	if err := models.DB.First(&shift, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "true",
			"message": "Shift not found",
		})
		return
	}

	if err := models.DB.Delete(&shift).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "true",
			"message": "Failed to delete shift",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   "false",
		"message": "Shift deleted successfully",
	})
}