package controllers

import (
	"errors"
	"net/http"

	"github.com/OrryFrasetyo/go-api-hotelqu/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// type validation department input
type ValidateDepartmentInput struct {
	ParentDepartmentId *int   `json:"parent_department_id"`
	DepartmentName     string `json:"department_name" binding:"required"`
}

// get all departments
func FindDepartments(c *gin.Context) {
	var departments []models.Department
	models.DB.Find(&departments)

	c.JSON(200, gin.H{
		"error": false,
		"message": "List Data Departments",
		"data":    departments,
	})
}

// store a department
func StoreDepartment(c *gin.Context) {
	var input ValidateDepartmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": true,
				"message": out,
			})
		}
		return
	}

	// create department
	department := models.Department{
		ParentDepartmentId: input.ParentDepartmentId,
		DepartmentName:     input.DepartmentName,
	}
	models.DB.Create(&department)

	// return response json
	c.JSON(201, gin.H{
		"error": false,
		"message": "Department Created Successfully",
		"data":    department,
	})
}

// get department by id
func FindDepartmentById(c *gin.Context) {
	var department models.Department
	if err := models.DB.Where("id = ?", c.Param("id")).First(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "Record not found!"})
		return
	}

	c.JSON(200, gin.H{
		"error": false,
		"message": "Detail Data Department by ID : " + c.Param("id"),
		"data":    department,
	})
}

// update department
func UpdateDepartment(c *gin.Context) {
	var department models.Department
	if err := models.DB.Where("id = ?", c.Param("id")).First(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "Record not found!",
		})
		return
	}

	var input ValidateDepartmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), GetErrorMsg(fe)}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": true,
				"message": out,
			})
		}
		return
	}

	models.DB.Model(&department).Updates(input)

	c.JSON(200, gin.H{
		"error": false,
		"message": "Department Updated Successfully",
		"data":    department,
	})
}

// delete department
func DeleteDepartment(c *gin.Context) {
	var department models.Department
	if err := models.DB.Where("id = ?", c.Param("id")).First(&department).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": "Record not found!"})
		return
	}

	models.DB.Delete(&department)

	c.JSON(200, gin.H{
		"error": false,
		"message": "Department Deleted Successfully",
	})
}
