package department

type ValidateDepartmentInput struct {
	ParentDepartmentId *int   `json:"parent_department_id"`
	DepartmentName     string `json:"department_name" binding:"required"`
}