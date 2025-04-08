package main

import (
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/attendance"
	"github.com/OrryFrasetyo/go-api-hotelqu/middlewares"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func main()  {
	// initialization Gin
	router := gin.Default()

	// Static file server for uploads
	router.Static("/uploads", "./uploads")

	// call database connection
	models.ConnectDatabase()

	// create route with method GET
	router.GET("/", func(c *gin.Context) {
		// return response JSON
		c.JSON(200, gin.H {
			"message": "Hello Hotelqu",
		})
	})

	// Auth routes
	router.POST("/api/register", controllers.Register)
	router.POST("/api/login", controllers.Login)

	// Protected routes (requiring JWT authentication)
	protected := router.Group("/api")
	protected.Use(middlewares.JWTAuth())
	{
		// User profile route
		protected.GET("/user", controllers.GetProfile)
		protected.PUT("/user", controllers.UpdateProfile)
		protected.GET("/schedules/department", controllers.ListDepartmentSchedules)
		protected.POST("/schedules", controllers.CreateSchedule)
		protected.PUT("/schedules/:id", controllers.UpdateSchedule)
		protected.DELETE("/schedules/:id", controllers.DeleteSchedule)

		// attendance endpoints
		protected.POST("/attendance", controllers.CreateAttendance)
		protected.PUT("/attendance", controllers.UpdateAttendance)
		protected.GET("/attendance", attendance.GetAttendanceLastThreeDays)
		protected.GET("/attendance/today", attendance.GetAttendanceToday)
		protected.GET("/attendance/month", attendance.GetAttendanceThisMonth)
	}

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

	// Shift routes
	router.GET("/api/shifts", controllers.FindShifts)
	router.POST("/api/shifts", controllers.StoreShift)
	router.GET("/api/shifts/:id", controllers.FindShiftById)
	router.PUT("/api/shifts/:id", controllers.UpdateShift)
	router.DELETE("/api/shifts/:id", controllers.DeleteShift)

	// start server with port 3000
	router.Run(":3000")
}