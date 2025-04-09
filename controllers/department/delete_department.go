package department

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func DeleteDepartment(c *gin.Context) {
	var department models.Department
	if err := models.DB.Where("id = ?", c.Param("id")).First(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Record not found!"})
		return
	}

	models.DB.Delete(&department)

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Department Deleted Successfully",
	})
}