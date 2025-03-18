package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// GetProfile returns the profile of the authenticated employee
func GetProfile(c *gin.Context) {
	// Get employee ID from the JWT token (set by middleware)
	employeeId, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Find employee by ID with position and department information
	var employee models.Employee
	if err := models.DB.Preload("Position.Department").First(&employee, employeeId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   true,
			"message": "Employee not found",
		})
		return
	}

	// Create response
	var photoURL *string
	if employee.Photo != nil {
		photoURL = employee.Photo
	}

	// Prepare department name
	var departmentName string
	if employee.Position.Department.DepartmentName != "" {
		departmentName = employee.Position.Department.DepartmentName
	}

	// Create response
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Profile retrieved successfully",
		"profile": gin.H{
			"id":         employee.Id,
			"name":       employee.Name,
			"email":      employee.Email,
			"phone":      employee.Phone,
			"position":   employee.Position.PositionName,
			"department": departmentName,
			"photo":      photoURL,
		},
	})
}

// UpdateProfileInput defines the input structure for updating profile
type UpdateProfileInput struct {
	Name     string `form:"name" binding:"required"`
	Password string `form:"password"`
	Phone    string `form:"phone" binding:"required"`
}

// UpdateProfile handles the profile update of the authenticated employee
func UpdateProfile(c *gin.Context) {
	// Get employee ID from the JWT token (set by middleware)
	employeeId, exists := c.Get("employeeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Unauthorized access",
		})
		return
	}

	// Find employee by ID
	var employee models.Employee
	if err := models.DB.First(&employee, employeeId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   true,
				"message": "Employee not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Database error",
			})
		}
		return
	}

	// Bind form data
	var input UpdateProfileInput
	if err := c.ShouldBind(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": out,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}
		return
	}

	// Handle file upload if a photo is provided
	file, err := c.FormFile("photo")
	if err != nil && err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid file upload",
		})
		return
	}

	var photoPath *string
	if file != nil {
		// Validate file size (max 2MB)
		if file.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "File too large (max 2MB)",
			})
			return
		}

		// Validate file type
		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Only JPG, JPEG, and PNG files are allowed",
			})
			return
		}

		// Create uploads directory if it doesn't exist
		uploadDir := "./uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err := os.MkdirAll(uploadDir, 0755)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   true,
					"message": "Failed to create upload directory",
				})
				return
			}
		}

		// Generate unique filename
		filename := fmt.Sprintf("%d_%d%s", employeeId, time.Now().Unix(), fileExt)
		filepath := filepath.Join(uploadDir, filename)

		// Save the file
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to open uploaded file",
			})
			return
		}
		defer src.Close()

		dst, err := os.Create(filepath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to create destination file",
			})
			return
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to save file",
			})
			return
		}

		// Set photo path for database (relative URL)
		relativePath := "/uploads/" + filename
		photoPath = &relativePath

		// Delete old photo if exists
		if employee.Photo != nil && *employee.Photo != "" {
			oldPhotoPath := "." + *employee.Photo
			if _, err := os.Stat(oldPhotoPath); err == nil {
				os.Remove(oldPhotoPath)
			}
		}
	}

	// Update employee data
	employee.Name = input.Name
	employee.Phone = input.Phone
	
	// Only update password if provided
	if input.Password != "" {
		if err := employee.HashPassword(input.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Failed to hash password",
			})
			return
		}
	}

	// Update photo if a new one was uploaded
	if photoPath != nil {
		employee.Photo = photoPath
	}

	// Save changes to database
	if err := models.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to update profile",
		})
		return
	}

	// Get updated employee with position and department
	if err := models.DB.Preload("Position.Department").First(&employee, employeeId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to retrieve updated profile",
		})
		return
	}

	// Prepare department name
	var departmentName string
	if employee.Position.Department.DepartmentName != "" {
		departmentName = employee.Position.Department.DepartmentName
	}

	// Return success response with updated profile
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Profile updated successfully",
		"profile": gin.H{
			"id":         employee.Id,
			"name":       employee.Name,
			"email":      employee.Email,
			"phone":      employee.Phone,
			"position":   employee.Position.PositionName,
			"department": departmentName,
			"photo":      employee.Photo,
		},
	})
}