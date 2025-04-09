package shift

import (
	"net/http"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func StoreShift(c *gin.Context) {
	var input ShiftInput

	if err := c.ShouldBindJSON(&input); err != nil {
		var errors []errormessage.ErrorMsg
		for _, e := range err.(validator.ValidationErrors) {
			errorMsg := errormessage.ErrorMsg{
				Field:   e.Field(),
				Message: errormessage.GetErrorMsg(e),
			}
			errors = append(errors, errorMsg)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "true",
			"message": "Validation failed",
			"data":    errors,
		})
		return
	}

	shift := models.Shift{
		Type:      input.Type,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
	}

	if err := models.DB.Create(&shift).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "true",
			"message": "Failed to create shift",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   "false",
		"message": "Shift created successfully",
		"data":    shift,
	})
}