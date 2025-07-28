package task


import (
	"net/http"
	"strconv"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// RestoreTask handles the restoration of a soft-deleted task by ID.
// Only managers/supervisors can restore tasks.
func RestoreTask(c *gin.Context) {
	// Get task ID from URL parameter
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "ID tugas tidak valid",
		})
		return
	}

	// Find the soft-deleted task using Unscoped
	var task models.Task
	if err := models.DB.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", taskID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Tugas yang dihapus tidak ditemukan",
		})
		return
	}

	// Authorization check
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

	// Load the task's employee with Position and Department
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
			"message": "Anda tidak memiliki izin untuk mengembalikan tugas di departemen lain",
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

	// Restore the task using Unscoped and Update deleted_at to NULL
	if err := tx.Unscoped().Model(&task).Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengembalikan tugas",
		})
		return
	}

	// Restore associated TaskItems
	if err := tx.Unscoped().Model(&models.TaskItem{}).Where("task_id = ?", task.ID).Update("deleted_at", nil).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengembalikan item tugas terkait",
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
		"message": "Tugas berhasil dikembalikan",
	})
}