package employee

import (
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// GetEmployeesForDepartment mengembalikan daftar nama karyawan berdasarkan departemen
func GetEmployeesForDepartment(c *gin.Context) {
	// Mendapatkan ID karyawan dari token JWT (diatur oleh middleware)
	employeeID, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Akses tidak diizinkan",
		})
		return
	}

	// Mendapatkan informasi karyawan yang login
	var employee models.Employee
	if err := models.DB.Preload("Position.Department").First(&employee, employeeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Karyawan tidak ditemukan",
		})
		return
	}

	// Mendapatkan departemen ID dari karyawan yang login
	departmentID := employee.Position.DepartmentId

	// Mendapatkan semua karyawan dari departemen yang sama
	var employees []models.Employee
	if err := models.DB.
		Joins("JOIN positions ON employees.position_id = positions.id").
		Where("positions.department_id = ?", departmentID).
		Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Gagal mengambil data karyawan: " + err.Error(),
		})
		return
	}

	// Menyiapkan response
	var employeeNames []gin.H
	for _, emp := range employees {
		employeeNames = append(employeeNames, gin.H{
			"employee_id": emp.Id,
			"name": emp.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"error":     false,
		"message":   "Data karyawan berhasil diambil",
		"employees": employeeNames,
	})
}
