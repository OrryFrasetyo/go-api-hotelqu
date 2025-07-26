package task

import (
	"errors"
	"net/http"
	"time"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CreateTask is handler for create task
func CreateTask(c *gin.Context) {
	// Mendapatkan ID karyawan dari token JWT (diatur oleh middleware)
	creatorID, _ := c.Get("employeeId")

	// Binding input dari request
	var input CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			errorMessages := make([]gin.H, len(ve))
			for i, fe := range ve {
				errorMessages[i] = gin.H{
					"field":   fe.Field(),
					"message": errormessage.GetErrorMsg(fe),
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Validasi gagal",
				"errors":  errorMessages,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	// Memeriksa apakah karyawan yang ditugaskan ada
	var employee models.Employee
	if err := models.DB.First(&employee, input.EmployeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Karyawan tidak ditemukan",
		})
		return
	}

	// Parse tanggal tugas dari format DD-MM-YYYY ke time.Time
	taskDate, err := time.Parse("02-01-2006", input.DateTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Format tanggal tugas tidak valid. Gunakan DD-MM-YYYY",
		})
		return
	}

	// Konversi ke format YYYY-MM-DD untuk database
	input.DateTask = taskDate.Format("2006-01-02")

	// Parse tanggal deadline dari format DD-MM-YYYY ke time.Time
	deadlineDate, err := time.Parse("02-01-2006", input.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Format tanggal deadline tidak valid. Gunakan DD-MM-YYYY",
		})
		return
	}

	// Konversi ke format YYYY-MM-DD untuk database
	input.Deadline = deadlineDate.Format("2006-01-02")

	// Validasi: Cek apakah ada jadwal kerja pada tanggal tersebut untuk employee yang ditugaskan
	var schedule models.Schedule
	if err := models.DB.Where("employee_id = ? AND date_schedule = ?", input.EmployeeID, input.DateTask).First(&schedule).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Tidak ada jadwal kerja pada tanggal tersebut",
		})
		return
	}

	// Membuat task baru (tanpa ScheduleID)
	task := models.Task{
		EmployeeID: input.EmployeeID,
		CreatedBy:  uint(creatorID.(int)),
		DateTask:   input.DateTask,
		Deadline:   input.Deadline,
		Status:     "Belum Dikerjakan",
		Message:    "-",
		Feedback:   "-",
	}

	// Menyimpan task ke database
	if err := models.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal membuat tugas: " + err.Error(),
		})
		return
	}

	// Membuat task items
	var taskItems []models.TaskItem
	for _, description := range input.TaskItems {
		taskItems = append(taskItems, models.TaskItem{
			TaskID:      task.ID,
			Description: description,
			IsCompleted: false,
		})
	}

	// Menyimpan task items ke database
	if err := models.DB.Create(&taskItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal membuat item tugas: " + err.Error(),
		})
		return
	}

	// Mengambil data task lengkap untuk response
	var createdTask models.Task
	models.DB.Preload("Employee").Preload("Creator").Preload("TaskItems").First(&createdTask, task.ID)

	// Menyiapkan response (tanpa schedule)
	response := gin.H{
		"error":   false,
		"message": "Tugas berhasil ditambahkan",
		"task": gin.H{
			"id": createdTask.ID,
			"employee": gin.H{
				"id":   createdTask.Employee.Id,
				"name": createdTask.Employee.Name,
			},
			"created_by": gin.H{
				"id":   createdTask.Creator.Id,
				"name": createdTask.Creator.Name,
			},
			"task_items": func() []gin.H {
				items := make([]gin.H, len(createdTask.TaskItems))
				for i, item := range createdTask.TaskItems {
					items[i] = gin.H{
						"id":           item.ID,
						"description":  item.Description,
						"is_completed": item.IsCompleted,
					}
				}
				return items
			}(),
			"date_task":  createdTask.DateTask,
			"deadline":   createdTask.Deadline,
			"status":     createdTask.Status,
			"feedback":   createdTask.Feedback,
			"message":    createdTask.Message,
			"created_at": createdTask.CreatedAt,
			"updated_at": createdTask.UpdatedAt,
		},
	}

	c.JSON(http.StatusCreated, response)
}
