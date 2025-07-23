package middlewares


import (
	"net/http"
	"strings"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
)

// manager auth is middleware for validation when employee is manajer or supervisor or etc
func ManagerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan ID karyawan dari token JWT (diatur oleh middleware JWTAuth)
		employeeID, exists := c.Get("employeeId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Akses tidak diizinkan",
			})
			c.Abort()
			return
		}

		// Mendapatkan informasi karyawan termasuk posisi
		var employee models.Employee
		if err := models.DB.Preload("Position").First(&employee, employeeID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "Karyawan tidak ditemukan",
			})
			c.Abort()
			return
		}

		// periksa apakah posisi karyawan mengandung kata "manajer" atau "supervisor" atau dll (case insensitive)
		positionName := strings.ToLower(employee.Position.PositionName)
		if !strings.Contains(positionName, "manager") && !strings.Contains(positionName, "supervisor") && !strings.Contains(positionName, "chief") && !strings.Contains(positionName, "executive") && !strings.Contains(positionName, "director") && !strings.Contains(positionName, "sous") && !strings.Contains(positionName, "partie") {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "Hanya manajer atau supervisor yang dapat mengakses fitur ini",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}