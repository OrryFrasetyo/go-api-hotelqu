package task

import (
	"net/http"
	"strconv"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// ChecklistTask handles task checklist updates by employees
func ChecklistTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "ID tugas tidak valid",
		})
		return
	}

	// Get employee ID from JWT token
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Akses tidak diizinkan",
		})
		return
	}

	// Bind request body
	var input ChecklistTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Data input tidak valid: " + err.Error(),
		})
		return
	}

	// Find the task and verify ownership
	var task models.Task
	if err := models.DB.Preload("Employee").Preload("Creator").Preload("Creator.Position").Preload("TaskItems").First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Tugas tidak ditemukan",
		})
		return
	}

	// Verify that the task belongs to the authenticated employee
	employeeIDInt, ok := employeeID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Error dalam mengambil ID karyawan",
		})
		return
	}

	if task.EmployeeID != uint(employeeIDInt) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "Anda tidak memiliki izin untuk mengupdate tugas ini",
		})
		return
	}

	// Start transaction
	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update task items
	for _, inputItem := range input.TaskItems {
		// Find the task item
		var taskItem models.TaskItem
		if err := tx.Where("id = ? AND task_id = ?", inputItem.ID, taskID).First(&taskItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "Item tugas dengan ID " + strconv.Itoa(int(inputItem.ID)) + " tidak ditemukan",
			})
			return
		}

		// Update is_completed status
		if inputItem.IsCompleted != nil {
			taskItem.IsCompleted = *inputItem.IsCompleted
		} else {
			// Default to false if not provided
			taskItem.IsCompleted = false
		}

		if err := tx.Save(&taskItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Gagal mengupdate item tugas",
			})
			return
		}
	}

	// Update task status and message
	updateData := map[string]interface{}{
		"status":   "Sedang Dicek",
		"feedback": "-",
	}
	if input.Message != nil {
		updateData["message"] = *input.Message
	}

	if err := tx.Model(&task).Updates(updateData).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengupdate tugas",
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menyimpan perubahan",
		})
		return
	}

	// Reload task with updated data
	if err := models.DB.Preload("Employee").Preload("Creator").Preload("Creator.Position").Preload("TaskItems").First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil data tugas yang diupdate",
		})
		return
	}

	// Format task items for response
	var taskItems []gin.H
	for _, item := range task.TaskItems {
		taskItems = append(taskItems, gin.H{
			"id":           item.ID,
			"description":  item.Description,
			"is_completed": item.IsCompleted,
		})
	}

	// Format dates properly
	var formattedDateTask, formattedDeadline string
	if task.DateTask != "" && task.DateTask != "0000-00-00" {
		if parsedDateTask, err := time.Parse("2006-01-02", task.DateTask); err == nil {
			formattedDateTask = parsedDateTask.Format("2006-01-02")
		} else {
			formattedDateTask = task.DateTask
		}
	} else {
		formattedDateTask = task.DateTask
	}

	if task.Deadline != "" && task.Deadline != "0000-00-00" {
		if parsedDeadline, err := time.Parse("2006-01-02", task.Deadline); err == nil {
			formattedDeadline = parsedDeadline.Format("2006-01-02")
		} else {
			formattedDeadline = task.Deadline
		}
	} else {
		formattedDeadline = task.Deadline
	}

	// Format response
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
		"message": "Tugas berhasil dikerjakan",
		"task":    formattedTask,
	})
}
