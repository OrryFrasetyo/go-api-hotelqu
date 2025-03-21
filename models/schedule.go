package models

import "time"

type Schedule struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	EmployeeID   uint      `json:"employee_id" gorm:"index"`
	Employee     Employee  `json:"employee" gorm:"foreignKey:EmployeeID"`
	ShiftID      uint      `json:"shift_id" gorm:"index"`
	Shift        Shift     `json:"shift" gorm:"foreignKey:ShiftID"`
	CreatedBy    uint      `json:"created_by" gorm:"index"`
	Creator      Employee  `json:"creator" gorm:"foreignKey:CreatedBy"`
	DateSchedule time.Time `json:"date_schedule" gorm:"type:date;index"`
	Status       string    `json:"status" gorm:"type:varchar(20);default:'hadir'"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
