package department

import (
	"errors"
	"net/http"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func StoreDepartment(c *gin.Context) {
	var input ValidateDepartmentInput
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

	department := models.Department{
		ParentDepartmentId: input.ParentDepartmentId,
		DepartmentName:     input.DepartmentName,
	}
	models.DB.Create(&department)

	c.JSON(201, gin.H{
		"error":   false,
		"message": "Department Created Successfully",
		"data":    department,
	})
}