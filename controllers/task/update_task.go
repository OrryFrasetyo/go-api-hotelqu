package task

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	errormessage "github.com/OrryFrasetyo/go-api-hotelqu/controllers/error_message"
	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UpdateTask menangani PUT /api/tasks/:id
func UpdateTask(c *gin.Context) {
	// Mendapatkan ID karyawan dari token JWT (diatur oleh middleware)
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Akses tidak diizinkan",
		})
		return
	}

	// Mendapatkan informasi karyawan termasuk posisi dan departemen
	var employee models.Employee
	if err := models.DB.Preload("Position").Preload("Position.Department").First(&employee, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Karyawan tidak ditemukan",
		})
		return
	}

	// Mendapatkan ID tugas dari parameter URL
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "ID tugas diperlukan",
		})
		return
	}

	// Konversi ID tugas ke uint
	taskIDUint, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "ID tugas tidak valid",
		})
		return
	}

	// Binding input dari request
	var input UpdateTaskInput
	if bindErr := c.ShouldBindJSON(&input); bindErr != nil {
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
			"message": bindErr.Error(),
		})
		return
	}

	// Mendapatkan tugas yang akan diupdate
	var task models.Task
	if dbErr := models.DB.Preload("Employee").Preload("Employee.Position").Preload("Creator").Preload("TaskItems").First(&task, taskIDUint).Error; dbErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Tugas tidak ditemukan",
		})
		return
	}

	// Memeriksa apakah karyawan yang mengedit adalah manajer/supervisor dari departemen yang sama dengan karyawan yang ditugaskan
	if task.Employee.Position.DepartmentId != employee.Position.DepartmentId {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "Anda hanya dapat mengedit tugas untuk karyawan di departemen Anda",
		})
		return
	}

	// Memulai transaksi database
	tx := models.DB.Begin()

	// Update employee_id (wajib)
	var newEmployee models.Employee
	if dbErr := tx.Preload("Position").First(&newEmployee, input.EmployeeID).Error; dbErr != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Karyawan yang ditugaskan tidak ditemukan",
		})
		return
	}

	// Memeriksa apakah karyawan yang ditugaskan berada di departemen yang sama
	if newEmployee.Position.DepartmentId != employee.Position.DepartmentId {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "Anda hanya dapat menugaskan karyawan di departemen Anda",
		})
		return
	}

	task.EmployeeID = input.EmployeeID

	// Update date_task (wajib)
	taskDate, err := time.Parse("02-01-2006", input.DateTask)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Format tanggal tugas tidak valid. Gunakan DD-MM-YYYY",
		})
		return
	}
	task.DateTask = taskDate.Format("2006-01-02")

	// Update deadline (wajib)
	deadlineDate, err := time.Parse("02-01-2006", input.Deadline)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Format tanggal deadline tidak valid. Gunakan DD-MM-YYYY",
		})
		return
	}
	task.Deadline = deadlineDate.Format("2006-01-02")

	// Update status (wajib)
	task.Status = input.Status

	// Update feedback jika disediakan
	if input.Feedback != nil {
		task.Feedback = *input.Feedback
	}

	// Validasi jadwal - cek apakah ada jadwal untuk employee_id pada date_task
	var schedule models.Schedule
	if err := tx.Where("employee_id = ? AND date_schedule = ?", task.EmployeeID, task.DateTask).First(&schedule).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Tidak ada jadwal kerja pada tanggal tersebut",
		})
		return
	}

	// Simpan perubahan pada task
	if err := tx.Save(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengupdate tugas: " + err.Error(),
		})
		return
	}

	// Update task items jika disediakan
	if len(input.TaskItems) > 0 {
		// Membuat map untuk menyimpan ID task item yang ada
		existingTaskItems := make(map[uint]bool)
		for _, item := range task.TaskItems {
			existingTaskItems[item.ID] = true
		}

		// Membuat map untuk menyimpan ID task item yang akan dipertahankan
		requestTaskItemIDs := make(map[uint]bool)

		// Memproses setiap task item
		for _, itemInput := range input.TaskItems {
			// Jika ID tidak disediakan, ini adalah item baru
			if itemInput.ID == nil {
				// Pastikan description disediakan untuk item baru
				if itemInput.Description == nil {
					tx.Rollback()
					c.JSON(http.StatusBadRequest, gin.H{
						"error":   true,
						"message": "Deskripsi diperlukan untuk item tugas baru",
					})
					return
				}

				// Nilai default untuk IsCompleted
				isCompleted := false
				if itemInput.IsCompleted != nil {
					isCompleted = *itemInput.IsCompleted
				}

				// Membuat task item baru
				newItem := models.TaskItem{
					TaskID:      task.ID,
					Description: *itemInput.Description,
					IsCompleted: isCompleted,
				}

				if err := tx.Create(&newItem).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   true,
						"message": "Gagal membuat item tugas baru: " + err.Error(),
					})
					return
				}
			} else {
				// Ini adalah update untuk item yang sudah ada
				// Tambahkan ID ke map item yang akan dipertahankan
				requestTaskItemIDs[*itemInput.ID] = true

				// Periksa apakah item ada
				if !existingTaskItems[*itemInput.ID] {
					tx.Rollback()
					c.JSON(http.StatusNotFound, gin.H{
						"error":   true,
						"message": "Item tugas tidak ditemukan",
					})
					return
				}

				// Mendapatkan item yang akan diupdate
				var taskItem models.TaskItem
				if err := tx.First(&taskItem, *itemInput.ID).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusNotFound, gin.H{
						"error":   true,
						"message": "Item tugas tidak ditemukan",
					})
					return
				}

				// Update description jika disediakan
				if itemInput.Description != nil {
					taskItem.Description = *itemInput.Description
				}

				// Update isCompleted jika disediakan
				if itemInput.IsCompleted != nil {
					taskItem.IsCompleted = *itemInput.IsCompleted
				}

				// Simpan perubahan pada task item
				if err := tx.Save(&taskItem).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   true,
						"message": "Gagal mengupdate item tugas: " + err.Error(),
					})
					return
				}
			}
		}

		// Hapus task items yang tidak ada dalam request
		for _, existingItem := range task.TaskItems {
			// Jika item tidak ada dalam request, hapus dari database
			if !requestTaskItemIDs[existingItem.ID] {
				if err := tx.Delete(&models.TaskItem{}, existingItem.ID).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   true,
						"message": "Gagal menghapus item tugas: " + err.Error(),
					})
					return
				}
			}
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menyimpan perubahan: " + err.Error(),
		})
		return
	}

	// Mengambil data task lengkap untuk response
	var updatedTask models.Task
	models.DB.Preload("Employee").Preload("Creator").Preload("TaskItems").First(&updatedTask, task.ID)

	// Konversi format tanggal dari YYYY-MM-DD ke DD-MM-YYYY untuk response
	dateTask, _ := time.Parse("2006-01-02", updatedTask.DateTask)
	deadline, _ := time.Parse("2006-01-02", updatedTask.Deadline)
	dateTaskFormatted := dateTask.Format("02-01-2006")
	deadlineFormatted := deadline.Format("02-01-2006")

	// Menyiapkan response
	response := gin.H{
		"error":   false,
		"message": "Tugas berhasil diedit",
		"task": gin.H{
			"id": updatedTask.ID,
			"employee": gin.H{
				"id":   updatedTask.Employee.Id,
				"name": updatedTask.Employee.Name,
			},
			"created_by": gin.H{
				"id":   updatedTask.Creator.Id,
				"name": updatedTask.Creator.Name,
			},
			"task_items": func() []gin.H {
				items := make([]gin.H, len(updatedTask.TaskItems))
				for i, item := range updatedTask.TaskItems {
					items[i] = gin.H{
						"id":           item.ID,
						"description":  item.Description,
						"is_completed": item.IsCompleted,
					}
				}
				return items
			}(),
			"date_task":  dateTaskFormatted,
			"deadline":   deadlineFormatted,
			"status":     updatedTask.Status,
			"message":    updatedTask.Message,
			"feedback":   updatedTask.Feedback,
			"created_at": updatedTask.CreatedAt,
			"updated_at": updatedTask.UpdatedAt,
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTaskStatus menangani PUT /api/tasks/status/:id
func UpdateTaskStatus(c *gin.Context) {

	// Mendapatkan ID karyawan dari token JWT (diatur oleh middleware)
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Akses tidak diizinkan",
		})
		return
	}

	// Mendapatkan informasi karyawan termasuk posisi
	var employee models.Employee
	if err := models.DB.Preload("Position").First(&employee, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Karyawan tidak ditemukan",
		})
		return
	}

	// Mendapatkan ID tugas dari parameter URL
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "ID tugas diperlukan",
		})
		return
	}

	// Konversi ID tugas ke uint
	taskIDUint, err := strconv.ParseUint(taskID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "ID tugas tidak valid",
		})
		return
	}

	// Binding input dari request
	var input UpdateTaskStatusInput
	if bindErr := c.ShouldBindJSON(&input); bindErr != nil {
		var ve validator.ValidationErrors
		if errors.As(bindErr, &ve) {
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
			"message": bindErr.Error(),
		})
		return
	}

	// Mendapatkan tugas yang akan diupdate
	var task models.Task
	if dbErr := models.DB.Preload("Employee").Preload("Employee.Position").First(&task, taskIDUint).Error; dbErr != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Tugas tidak ditemukan",
		})
		return
	}

	// Memeriksa apakah karyawan yang mengedit adalah manajer/supervisor dari departemen yang sama
	if task.Employee.Position.DepartmentId != employee.Position.DepartmentId {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "Anda hanya dapat mengedit tugas untuk karyawan di departemen Anda",
		})
		return
	}

	// Update status (wajib)
	task.Status = input.Status

	// Update feedback jika disediakan
	if input.Feedback != nil {
		task.Feedback = *input.Feedback
	}

	// Simpan perubahan pada task - hanya update field yang diperlukan
	updateData := map[string]interface{}{
		"status": task.Status,
	}
	if input.Feedback != nil {
		updateData["feedback"] = task.Feedback
	}

	if err := models.DB.Model(&task).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengupdate status tugas: " + err.Error(),
		})
		return
	}

	// Response sukses
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Status tugas dan feedback berhasil diedit",
	})
}
