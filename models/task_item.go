package models

import "time"

type TaskItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TaskID      uint      `json:"task_id" gorm:"index"`
	Description string    `json:"description" gorm:"type:text"`
	IsCompleted bool      `json:"is_completed" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}