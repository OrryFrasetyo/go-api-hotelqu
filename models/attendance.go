package models

import "time"

type Attendance struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	ScheduleID     uint      `json:"schedule_id" gorm:"index"`
	Schedule       Schedule  `json:"schedule" gorm:"foreignKey:ScheduleID"`
	Date           string    `json:"date" gorm:"type:date"`
	ClockIn        string    `json:"clock_in" gorm:"type:varchar(8)"`
	ClockOut       string    `json:"clock_out" gorm:"type:varchar(8)"`
	Duration       string    `json:"duration" gorm:"type:varchar(30)"`
	ClockInStatus  string    `json:"clock_in_status" gorm:"type:varchar(20)"`
	ClockOutStatus string    `json:"clock_out_status" gorm:"type:varchar(20)"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}