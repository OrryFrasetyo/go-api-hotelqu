package models


type Shift struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type" gorm:"type:varchar(30);not null"`
	StartTime string    `json:"start_time" gorm:"type:varchar(8);not null"`
	EndTime   string    `json:"end_time" gorm:"type:varchar(8);not null"`
}