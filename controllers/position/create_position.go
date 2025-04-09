package position

import (
	"errors"
	"net/http"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func StorePosition(c *gin.Context) {
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

	position := models.Position{
		DepartmentId: input.DepartmentId,
		PositionName: input.PositionName,
		IsCompleted:  input.IsCompleted,
	}
	models.DB.Create(&position)

	c.JSON(201, gin.H{
		"error":   false,
		"message": "Position Created Successfully",
		"data": map[string]interface{}{
			"id":              position.Id,
			"department_id":   position.DepartmentId,
			"department_name": department.DepartmentName,
			"position_name":   position.PositionName,
			"is_completed":    position.IsCompleted,
		},
	})
}