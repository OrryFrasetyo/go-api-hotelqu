package main

import (
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/attendance"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/authentication"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/department"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/employee"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/position"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/schedule"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/shift"
	"github.com/OrryFrasetyo/go-api-hotelqu/middlewares"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func main()  {
	router := gin.Default()

	router.Static("/uploads", "./uploads")

	models.ConnectDatabase()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H {
			"message": "Hello Hotelqu",
		})
	})

	// Auth routes
	router.POST("/api/register", authentication.Register)
	router.POST("/api/login", authentication.Login)

	// Protected routes (requiring JWT authentication)
	protected := router.Group("/api")
	protected.Use(middlewares.JWTAuth())
	{
		// User profile route
		protected.GET("/user", employee.GetProfile)
		protected.PUT("/user", employee.UpdateProfile)

		// schedules endpoints
		protected.GET("/schedules/department", schedule.ListDepartmentSchedules)
		protected.POST("/schedules", schedule.CreateSchedule)
		protected.PUT("/schedules/:id", schedule.UpdateSchedule)
		protected.DELETE("/schedules/:id", schedule.DeleteSchedule)

		// attendance endpoints
		protected.POST("/attendance", attendance.CreateAttendance)
		protected.PUT("/attendance", attendance.UpdateAttendance)
		protected.GET("/attendance", attendance.GetAttendanceLastThreeDays)
		protected.GET("/attendance/today", attendance.GetAttendanceToday)
		protected.GET("/attendance/month", attendance.GetAttendanceThisMonth)
		protected.GET("/attendance/status", attendance.GetAttendanceByStatus)
	}

	// Department routes
	router.GET("/api/departments", department.FindDepartments)
	router.POST("/api/departments", department.StoreDepartment)
	router.GET("/api/departments/:id", department.FindDepartmentById)
	router.PUT("/api/departments/:id", department.UpdateDepartment)
	router.DELETE("/api/departments/:id", department.DeleteDepartment)

	// Position routes
	router.GET("/api/positions", position.FindPositions)
	router.POST("/api/positions", position.StorePosition)
	router.GET("/api/positions/:id", position.FindPositionById)
	router.PUT("/api/positions/:id", position.UpdatePosition)
	router.DELETE("/api/positions/:id", position.DeletePosition)

	// Shift routes
	router.GET("/api/shifts", shift.FindShifts)
	router.POST("/api/shifts", shift.StoreShift)
	router.GET("/api/shifts/:id", shift.FindShiftById)
	router.PUT("/api/shifts/:id", shift.UpdateShift)
	router.DELETE("/api/shifts/:id", shift.DeleteShift)

	// start server with port 3000
	router.Run(":3000")
}