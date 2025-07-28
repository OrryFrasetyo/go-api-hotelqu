package task

import (
	"net/http"
	"strconv"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// DeleteTask handles the soft deletion of a task by ID.
// Only managers/supervisors can delete tasks.
func DeleteTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "ID tugas tidak valid",
		})
		return
	}

	// Find the task (exclude soft deleted)
	var task models.Task
	if err := models.DB.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Tugas tidak ditemukan",
		})
		return
	}

	// Check if the employee (manager/supervisor) is in the same department as the task's employee
	// Get employee ID from JWT token (set by JWTAuth middleware)
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Akses tidak diizinkan",
		})
		return
	}

	// Load manager with Position and Department
	var manager models.Employee
	if err := models.DB.Preload("Position.Department").First(&manager, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Karyawan tidak ditemukan",
		})
		return
	}

	// Load the task's employee with Position and Department to get their department ID
	var taskEmployee models.Employee
	if err := models.DB.Preload("Position.Department").First(&taskEmployee, task.EmployeeID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil data karyawan tugas",
		})
		return
	}

	// Check if the manager's department ID matches the task's employee's department ID
	if manager.Position.DepartmentId != taskEmployee.Position.DepartmentId {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   true,
			"message": "Anda tidak memiliki izin untuk menghapus tugas di departemen lain",
		})
		return
	}

	// Start database transaction
	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Soft delete associated TaskItems first
	if err := tx.Where("task_id = ?", task.ID).Delete(&models.TaskItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menghapus item tugas terkait",
		})
		return
	}

	// Soft delete the task
	if err := tx.Delete(&task).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal menghapus tugas",
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

	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Tugas berhasil dihapus (soft delete)",
	})
}
