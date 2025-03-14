package controllers

import (
	"errors"
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// type validation position input
type ValidatePositionInput struct {
	DepartmentId  int    `json:"department_id" binding:"required"`
	PositionName  string `json:"position_name" binding:"required"`
	IsCompleted   bool   `json:"is_completed"`
}

// get all positions
func FindPositions(c *gin.Context) {
	var positions []models.Position
	
	// Preload Department untuk mendapatkan department_name
	result := models.DB.Preload("Department").Find(&positions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"message": result.Error.Error(),
		})
		return
	}

	// Menyesuaikan format response sesuai API spec
	type ResponsePosition struct {
		Id            int    `json:"id"`
		DepartmentId  int    `json:"department_id"`
		DepartmentName string `json:"department_name"`
		PositionName  string `json:"position_name"`
		IsCompleted   bool   `json:"is_completed"`
	}

	var responsePositions []ResponsePosition
	for _, position := range positions {
		responsePosition := ResponsePosition{
			Id:            position.Id,
			DepartmentId:  position.DepartmentId,
			DepartmentName: position.Department.DepartmentName,
			PositionName:  position.PositionName,
			IsCompleted:   position.IsCompleted,
		}
		responsePositions = append(responsePositions, responsePosition)
	}

	c.JSON(200, gin.H{
		"error": false,
		"message": "List Data Positions",
		"data":    responsePositions,
	})
}

// store a position
func StorePosition(c *gin.Context) {
	var input ValidatePositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": true,
				"message": out,
			})
		}
		return
	}

	// verify if department_id valid?
	var department models.Department
	if err := models.DB.Where("id = ?", input.DepartmentId).First(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "Department not found!",
		})
		return
	}

	// create position
	position := models.Position{
		DepartmentId: input.DepartmentId,
		PositionName: input.PositionName,
		IsCompleted:  input.IsCompleted,
	}
	models.DB.Create(&position)

	// return response json
	c.JSON(201, gin.H{
		"error": false,
		"message": "Position Created Successfully",
		"data": map[string]interface{}{
			"id":             position.Id,
			"department_id":  position.DepartmentId,
			"department_name": department.DepartmentName,
			"position_name":  position.PositionName,
			"is_completed":   position.IsCompleted,
		},
	})
}

// get position by id
func FindPositionById(c *gin.Context) {
	var position models.Position
	if err := models.DB.Preload("Department").Where("id = ?", c.Param("id")).First(&position).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "Record not found!",
		})
		return
	}

	c.JSON(200, gin.H{
		"error": false,
		"message": "Detail Data Position by ID : " + c.Param("id"),
		"data": map[string]interface{}{
			"id":             position.Id,
			"department_id":  position.DepartmentId,
			"department_name": position.Department.DepartmentName,
			"position_name":  position.PositionName,
			"is_completed":   position.IsCompleted,
		},
	})
}

// update position
func UpdatePosition(c *gin.Context) {
	var position models.Position
	if err := models.DB.Where("id = ?", c.Param("id")).First(&position).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "Record not found!",
		})
		return
	}

	var input ValidatePositionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": true,
				"message": out,
			})
		}
		return
	}

	// verify if department_id valid?
	var department models.Department
	if err := models.DB.Where("id = ?", input.DepartmentId).First(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "Department not found!",
		})
		return
	}

	models.DB.Model(&position).Updates(input)

	// Fetch updated position with department
	models.DB.Preload("Department").Where("id = ?", position.Id).First(&position)

	c.JSON(200, gin.H{
		"error": false,
		"message": "Position Updated Successfully",
		"data": map[string]interface{}{
			"id":             position.Id,
			"department_id":  position.DepartmentId,
			"department_name": position.Department.DepartmentName,
			"position_name":  position.PositionName,
			"is_completed":   position.IsCompleted,
		},
	})
}

// delete position
func DeletePosition(c *gin.Context) {
	var position models.Position
	if err := models.DB.Where("id = ?", c.Param("id")).First(&position).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "Record not found!",
		})
		return
	}

	models.DB.Delete(&position)

	c.JSON(200, gin.H{
		"error": false,
		"message": "Position Deleted Successfully",
	})
}



