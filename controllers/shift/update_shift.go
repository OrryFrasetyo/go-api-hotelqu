package shift

import (
	"net/http"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func UpdateShift(c *gin.Context) {
	var shift models.Shift
	id := c.Param("id")

	if err := models.DB.First(&shift, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "true",
			"message": "Shift not found",
		})
		return
	}

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

	updateData := models.Shift{
		Type:      input.Type,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
	}

	if err := models.DB.Model(&shift).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "true",
			"message": "Failed to update shift",
		})
		return
	}

	models.DB.First(&shift, id)

	c.JSON(http.StatusOK, gin.H{
		"error":   "false",
		"message": "Shift updated successfully",
		"data":    shift,
	})
}