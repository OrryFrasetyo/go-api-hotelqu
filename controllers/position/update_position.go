package position

import (
	"errors"
	"net/http"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func UpdatePosition(c *gin.Context) {
	var position models.Position
	if err := models.DB.Where("id = ?", c.Param("id")).First(&position).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Record not found!",
		})
		return
	}

	var input ValidatePositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]errormessage.ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = errormessage.ErrorMsg{Field: fe.Field(), Message: errormessage.GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": out,
			})
		}
		return
	}

	var department models.Department
	if err := models.DB.Where("id = ?", input.DepartmentId).First(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Department not found!",
		})
		return
	}

	models.DB.Model(&position).Updates(input)

	// Fetch updated position with department
	models.DB.Preload("Department").Where("id = ?", position.Id).First(&position)

	c.JSON(200, gin.H{
		"error":   false,
		"message": "Position Updated Successfully",
		"data": map[string]interface{}{
			"id":              position.Id,
			"department_id":   position.DepartmentId,
			"department_name": position.Department.DepartmentName,
			"position_name":   position.PositionName,
			"is_completed":    position.IsCompleted,
		},
	})
}