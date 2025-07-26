package task

import (
	"net/http"
	"strconv"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// ListDepartmentTasks menangani GET /api/tasks/department
func ListDepartmentTasks(c *gin.Context) {
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

	// Mendapatkan departemen ID dari posisi karyawan atau dari parameter query
	departmentID := c.Query("department_id")
	var deptID int

	if departmentID != "" {
		id, err := strconv.Atoi(departmentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "ID departemen tidak valid",
			})
			return
		}
		deptID = id
	} else {
		// Menggunakan departemen dari posisi karyawan
		deptID = employee.Position.DepartmentId
	}

	// Mendapatkan parameter tanggal (opsional)
	dateParam := c.Query("date_task")
	var taskDate time.Time
	var err error

	if dateParam != "" {
		// Parse tanggal dalam format DD-MM-YYYY
		taskDate, err = time.Parse("02-01-2006", dateParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Format tanggal tidak valid. Gunakan DD-MM-YYYY",
			})
			return
		}
	} else {
		taskDate = time.Now()
	}

	// Format tanggal untuk query database (YYYY-MM-DD)
	formattedDate := taskDate.Format("2006-01-02")

	// Mendapatkan informasi departemen
	var department models.Department
	if err := models.DB.First(&department, deptID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Departemen tidak ditemukan",
		})
		return
	}

	// Menghitung total karyawan di departemen ini
	var totalEmployees int64
	models.DB.Model(&models.Employee{}).
		Joins("JOIN positions ON employees.position_id = positions.id").
		Where("positions.department_id = ?", deptID).
		Count(&totalEmployees)

	// Mendapatkan semua tugas untuk karyawan di departemen pada tanggal yang ditentukan
	var tasks []models.Task
	query := models.DB.Preload("Employee").Preload("Creator").Preload("TaskItems").
		Joins("JOIN employees ON tasks.employee_id = employees.id").
		Joins("JOIN positions ON employees.position_id = positions.id").
		Where("positions.department_id = ? AND tasks.date_task = ?", deptID, formattedDate)

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil data tugas: " + err.Error(),
		})
		return
	}

	// Menyiapkan response
	taskList := make([]gin.H, 0, len(tasks))
	for _, task := range tasks {
		// Konversi format tanggal dari YYYY-MM-DD ke DD-MM-YYYY untuk response
		dateTask, _ := time.Parse("2006-01-02", task.DateTask)
		deadline, _ := time.Parse("2006-01-02", task.Deadline)
		dateTaskFormatted := dateTask.Format("02-01-2006")
		deadlineFormatted := deadline.Format("02-01-2006")

		// Menyiapkan task items
		taskItems := make([]gin.H, 0, len(task.TaskItems))
		for _, item := range task.TaskItems {
			taskItems = append(taskItems, gin.H{
				"id":           item.ID,
				"description":  item.Description,
				"is_completed": item.IsCompleted,
			})
		}

		taskList = append(taskList, gin.H{
			"id": task.ID,
			"employee": gin.H{
				"id":   task.Employee.Id,
				"name": task.Employee.Name,
			},
			"created_by": gin.H{
				"id":   task.Creator.Id,
				"name": task.Creator.Name,
			},
			"task_items": taskItems,
			"date_task":  dateTaskFormatted,
			"deadline":   deadlineFormatted,
			"status":     task.Status,
			"feedback":   task.Feedback,
			"created_at": task.CreatedAt,
			"updated_at": task.UpdatedAt,
		})
	}

	// Mengembalikan response
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Tugas Pegawai Berhasil Ditampilkan",
		"meta": gin.H{
			"date": taskDate.Format("02-01-2006"),
			"department": gin.H{
				"id":              department.Id,
				"department_name": department.DepartmentName,
			},
			"total_employees": totalEmployees,
		},
		"list_task": taskList,
	})
}
