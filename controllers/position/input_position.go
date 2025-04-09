package position

type ValidatePositionInput struct {
	DepartmentId int    `json:"department_id" binding:"required"`
	PositionName string `json:"position_name" binding:"required"`
	IsCompleted  bool   `json:"is_completed"`
}