package controllers

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ShiftInput struct {
	Type      string `json:"type" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

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

func StoreShift(c *gin.Context) {
	var input ShiftInput

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		var errors []ErrorMsg
		for _, e := range err.(validator.ValidationErrors) {
			errorMsg := ErrorMsg{
				Field:   e.Field(),
				Message: GetErrorMsg(e),
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

	// Create shift
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

func FindShiftById(c *gin.Context) {
	// Get ID from URL parameter
	var shift models.Shift
	id := c.Param("id")

	if err := models.DB.First(&shift, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "true",
			"message": "Shift not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   "false",
		"message": "Shift found",
		"data":    shift,
	})
}

func UpdateShift(c *gin.Context) {
	// Get ID from URL parameter
	var shift models.Shift
	id := c.Param("id")

	// Check if shift exists
	if err := models.DB.First(&shift, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "true",
			"message": "Shift not found",
		})
		return
	}

	// Validate input
	var input ShiftInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var errors []ErrorMsg
		for _, e := range err.(validator.ValidationErrors) {
			errorMsg := ErrorMsg{
				Field:   e.Field(),
				Message: GetErrorMsg(e),
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

	// Update shift
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

	// Fetch updated shift
	models.DB.First(&shift, id)

	c.JSON(http.StatusOK, gin.H{
		"error":   "false",
		"message": "Shift updated successfully",
		"data":    shift,
	})
}

func DeleteShift(c *gin.Context) {
	// Get ID from URL parameter
	var shift models.Shift
	id := c.Param("id")

	// Check if shift exists
	if err := models.DB.First(&shift, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "true",
			"message": "Shift not found",
		})
		return
	}

	// Delete shift
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