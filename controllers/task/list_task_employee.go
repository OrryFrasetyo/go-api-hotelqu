package task

import (
	"net/http"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// ListTaskEmployee handles listing tasks for the authenticated employee
// Can be filtered by date_task parameter
func ListTaskEmployee(c *gin.Context) {
	// Get employee ID from JWT token (set by JWTAuth middleware)
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Akses tidak diizinkan",
		})
		return
	}

	// Get date_task parameter (optional)
	dateTask := c.Query("date_task")

	// Build query
	query := models.DB.Where("employee_id = ?", employeeID)

	// Add date filter if provided
	if dateTask != "" {
		// Parse date from DD-MM-YYYY to YYYY-MM-DD format
		parsedDate, err := time.Parse("02-01-2006", dateTask)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Format tanggal tidak valid. Gunakan format DD-MM-YYYY",
			})
			return
		}
		formattedDate := parsedDate.Format("2006-01-02")
		query = query.Where("date_task = ?", formattedDate)
	}

	// Find task with preloaded relations (limit to 1 since we expect only one task per date)
	var task models.Task
	if err := query.Preload("Employee").Preload("Creator").Preload("Creator.Position").Preload("TaskItems").Order("date_task DESC, created_at DESC").First(&task).Error; err != nil {
		message := "Tidak ada tugas ditemukan"
		if dateTask != "" {
			message = "Tidak ada tugas ditemukan untuk tanggal " + dateTask
		}
		c.JSON(http.StatusOK, gin.H{
			"error":   false,
			"message": message,
			"task":    nil,
		})
		return
	}

	// Format task items
	var taskItems []gin.H
	for _, item := range task.TaskItems {
		taskItems = append(taskItems, gin.H{
			"id":           item.ID,
			"description":  item.Description,
			"is_completed": item.IsCompleted,
		})
	}

	// Format dates properly - handle empty dates
	var formattedDateTask, formattedDeadline string
	if task.DateTask != "" && task.DateTask != "0000-00-00" {
		if parsedDateTask, err := time.Parse("2006-01-02", task.DateTask); err == nil {
			formattedDateTask = parsedDateTask.Format("2006-01-02")
		} else {
			formattedDateTask = task.DateTask // fallback to original value
		}
	} else {
		formattedDateTask = task.DateTask
	}

	if task.Deadline != "" && task.Deadline != "0000-00-00" {
		if parsedDeadline, err := time.Parse("2006-01-02", task.Deadline); err == nil {
			formattedDeadline = parsedDeadline.Format("2006-01-02")
		} else {
			formattedDeadline = task.Deadline // fallback to original value
		}
	} else {
		formattedDeadline = task.Deadline
	}

	// Format single task response
	formattedTask := gin.H{
		"id": task.ID,
		"employee": gin.H{
			"id":   task.Employee.Id,
			"name": task.Employee.Name,
		},
		"created_by": gin.H{
			"id":       task.Creator.Id,
			"name":     task.Creator.Name,
			"position": task.Creator.Position.PositionName,
		},
		"task_items": taskItems,
		"date_task":  formattedDateTask,
		"deadline":   formattedDeadline,
		"status":     task.Status,
		"message":    task.Message,
		"feedback":   task.Feedback,
		"created_at": task.CreatedAt,
		"updated_at": task.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Tugas Pegawai Berhasil Ditampilkan",
		"task":    formattedTask,
	})
}
