package position

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func FindPositions(c *gin.Context) {
	var positions []models.Position

	result := models.DB.Preload("Department").Find(&positions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": result.Error.Error(),
		})
		return
	}

	type ResponsePosition struct {
		Id             int    `json:"id"`
		DepartmentId   int    `json:"department_id"`
		DepartmentName string `json:"department_name"`
		PositionName   string `json:"position_name"`
		IsCompleted    bool   `json:"is_completed"`
	}

	var responsePositions []ResponsePosition
	for _, position := range positions {
		responsePosition := ResponsePosition{
			Id:             position.Id,
			DepartmentId:   position.DepartmentId,
			DepartmentName: position.Department.DepartmentName,
			PositionName:   position.PositionName,
			IsCompleted:    position.IsCompleted,
		}
		responsePositions = append(responsePositions, responsePosition)
	}

	c.JSON(200, gin.H{
		"error":   false,
		"message": "List Data Positions",
		"data":    responsePositions,
	})
}