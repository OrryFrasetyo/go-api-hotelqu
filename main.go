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

	// create route get all departments
	router.GET("/api/departments", controllers.FindDepartments)

	// create route store department
	router.POST("/api/departments", controllers.StoreDepartment)

	// create route detail department
	router.GET("/api/departments/:id", controllers.FindDepartmentById)

	// create route update department
	router.PUT("/api/departments/:id", controllers.UpdateDepartment)

	// create route delete department
	router.DELETE("/api/departments/:id", controllers.DeleteDepartment)

	// start server with port 3000
	router.Run(":3000")
}