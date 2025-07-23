package models


import "time"

type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee  Employee   `json:"employee" gorm:"foreignKey:EmployeeID"`
	ScheduleID uint      `json:"schedule_id" gorm:"index"`
	Schedule  Schedule   `json:"schedule" gorm:"foreignKey:ScheduleID"`
	CreatedBy uint      `json:"created_by" gorm:"index"`
	Creator   Employee   `json:"creator" gorm:"foreignKey:CreatedBy"`
	DateTask  string    `json:"date_task" gorm:"type:date"`
	Deadline  string    `json:"deadline" gorm:"type:date"`
	Status    string    `json:"status" gorm:"type:varchar(50);default:'Belum Dikerjakan'"`
	Message   string    `json:"message" gorm:"type:text"`
	Feedback  string    `json:"feedback" gorm:"type:text"`
	TaskItems []TaskItem `json:"task_items" gorm:"foreignKey:TaskID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}