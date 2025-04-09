package shift

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func FindShifts(c *gin.Context) {
	var shifts []models.Shift

	if err := models.DB.Find(&shifts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "true",
			"message": "Failed to fetch shifts",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   "false",
		"message": "Shifts fetched successfully",
		"data":    shifts,
	})
}