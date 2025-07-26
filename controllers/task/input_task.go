package task

type CreateTaskInput struct {
	EmployeeID uint     `json:"employee_id" binding:"required"`
	TaskItems  []string `json:"task_items" binding:"required,min=1"`
	DateTask   string   `json:"date_task" binding:"required"`
	Deadline   string   `json:"deadline" binding:"required"`
}

type UpdateTaskInput struct {
	EmployeeID uint             `json:"employee_id" binding:"required"`  
	TaskItems  []UpdateTaskItem `json:"task_items" binding:"required,min=1"` 
	DateTask   string           `json:"date_task" binding:"required"`    
	Deadline   string           `json:"deadline" binding:"required"`     
	Status     string           `json:"status" binding:"required"`       
	Feedback   *string          `json:"feedback"`                        
}

type UpdateTaskItem struct {
	ID          *uint   `json:"id"`
	Description *string `json:"description"`
	IsCompleted *bool   `json:"is_completed"`
}
