package main

import (
	"log"
	"os"

	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/attendance"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/authentication"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/department"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/employee"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/position"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/schedule"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/shift"
	"github.com/OrryFrasetyo/go-api-hotelqu/controllers/task" // Tambahkan import untuk task
	"github.com/OrryFrasetyo/go-api-hotelqu/middlewares"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Static("/uploads", "./uploads")

	models.ConnectDatabase()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
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
		protected.GET("/employees", employee.GetEmployeesForDepartment)

		// schedules endpoints
		protected.GET("/schedules", schedule.ListSchedules)
		protected.GET("/schedules/today", schedule.GetTodaySchedule)
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

		// task endpoints (hanya untuk manajer/supervisor)
		taskRoutes := protected.Group("/tasks")
		taskRoutes.Use(middlewares.ManagerAuth())
		{
			taskRoutes.POST("", task.CreateTask)
			taskRoutes.GET("/department", task.ListDepartmentTasks)
			taskRoutes.PUT("/:id", task.UpdateTask)
			taskRoutes.PUT("/status/:id", task.UpdateTaskStatus)
		}
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
	// router.Run(":3000")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Port default jika dijalankan di lokal
	}

	log.Printf("Server starting on port %s", port)
	router.Run(":" + port)
}
