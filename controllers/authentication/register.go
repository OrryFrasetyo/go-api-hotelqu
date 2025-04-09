package authentication

import (
	"errors"
	"net/http"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var input RegisterInput
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

	var existingEmployee models.Employee
	if err := models.DB.Where("email = ?", input.Email).First(&existingEmployee).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Email already registered",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Database error",
		})
		return
	}

	var position models.Position
	result := models.DB.Where("position_name = ?", input.Position).First(&position)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Position '" + input.Position + "' not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Database error when searching position: " + result.Error.Error(),
			})
		}
		return
	}

	// create new employee
	employee := models.Employee{
		Name:       input.Name,
		Email:      input.Email,
		Password:   input.Password,
		Phone:      input.Phone,
		PositionId: position.Id,
	}

	if err := models.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to register employee",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Registration successful",
	})
}