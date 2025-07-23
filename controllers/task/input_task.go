package task

type CreateTaskInput struct {
	EmployeeID uint     `json:"employee_id" binding:"required"`
	ScheduleID uint     `json:"schedule_id" binding:"required"`
	TaskItems  []string `json:"task_items" binding:"required,min=1"`
	DateTask   string   `json:"date_task" binding:"required"`
	Deadline   string   `json:"deadline" binding:"required"`
}