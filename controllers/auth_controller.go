package controllers

import (
	"errors"
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/OrryFrasetyo/go-api-hotelqu/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Phone    string `json:"phone" binding:"required"`
	Position string `json:"position" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
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

	// validate email is registered?
	var existingEmployee models.Employee
	if err := models.DB.Where("email = ?", input.Email).First(&existingEmployee).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Email already registered",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Database error",
		})
		return
	}

	// Search position by position name
	// var position models.Position
	// if err := models.DB.Where("position_name = ?", input.Position).First(&position).Error; err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error":   true,
	// 		"message": "Position not found",
	// 	})
	// 	return
	// }
	var position models.Position
    result := models.DB.Where("position_name = ?", input.Position).First(&position)
    if result.Error != nil {
        // Jika position tidak ditemukan, berikan pesan error yang lebih spesifik
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": true,
                "message": "Position '" + input.Position + "' not found",
            })
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": true,
                "message": "Database error when searching position: " + result.Error.Error(),
            })
        }
        return
    }

	// create new employee
	employee := models.Employee{
		Name:       input.Name,
		Email:      input.Email,
		Password:   input.Password, 
		Phone:      input.Phone,
		PositionId: position.Id,
	}

	if err := models.DB.Create(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to register employee",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Registration successful",
	})
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
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

	// search employee by email 
	var employee models.Employee
	if err := models.DB.Where("email = ?", input.Email).First(&employee).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Invalid email or password",
		})
		return
	}

	// verification password
	if err := employee.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Invalid email or password",
		})
		return
	}

	// Generate token JWT
	token, err := utils.GenerateToken(employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to generate token",
		})
		return
	}

	// Kirim respons dengan token
	c.JSON(http.StatusOK, gin.H{
		"error":   false,
		"message": "Login successful",
		"loginResult": gin.H{
			"id":    employee.Id,
			"name":  employee.Name,
			"token": token,
		},
	})
}
