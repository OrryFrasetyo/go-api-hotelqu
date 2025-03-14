package main

import (
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func main()  {
	// initialization Gin
	router := gin.Default()

	// call database connection
	models.ConnectDatabase()

	// create route with method GET
	router.GET("/", func(c *gin.Context) {
		// return response JSON
		c.JSON(200, gin.H {
			"message": "Hello Hotelqu",
		})
	})

	// Department routes
	router.GET("/api/departments", controllers.FindDepartments)
	router.POST("/api/departments", controllers.StoreDepartment)
	router.GET("/api/departments/:id", controllers.FindDepartmentById)
	router.PUT("/api/departments/:id", controllers.UpdateDepartment)
	router.DELETE("/api/departments/:id", controllers.DeleteDepartment)

	// Position routes
	router.GET("/api/positions", controllers.FindPositions)
	router.POST("/api/positions", controllers.StorePosition)
	router.GET("/api/positions/:id", controllers.FindPositionById)
	router.PUT("/api/positions/:id", controllers.UpdatePosition)
	router.DELETE("/api/positions/:id", controllers.DeletePosition)

	// start server with port 3000
	router.Run(":3000")
}