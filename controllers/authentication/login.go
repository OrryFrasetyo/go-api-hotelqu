package authentication

import (
	"errors"
	"net/http"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/OrryFrasetyo/go-api-hotelqu/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Login(c *gin.Context) {
	var input LoginInput
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
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}
		return
	}

	// search employee by email
	var employee models.Employee
	if err := models.DB.Where("email = ?", input.Email).First(&employee).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Invalid email or password",
		})
		return
	}

	// verification password
	if err := employee.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Invalid email or password",
		})
		return
	}

	// Generate token JWT
	token, err := utils.GenerateToken(employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Login successful",
		"loginResult": gin.H{
			"id":    employee.Id,
			"name":  employee.Name,
			"token": token,
		},
	})
}