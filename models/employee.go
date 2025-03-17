package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Employee struct {
	Id         int      `json:"id" gorm:"primaryKey"`
	PositionId int      `json:"position_id" gorm:"index"`
	Position   Position `json:"-" gorm:"foreignKey:PositionId"`
	Name       string   `json:"name" gorm:"type:varchar(100);not null"`
	Email      string   `json:"email" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password   string   `json:"-" gorm:"type:varchar(255);not null"`
	Photo      *string  `json:"photo" gorm:"type:varchar(100)"`
	Phone      string   `json:"phone" gorm:"type:varchar(16);not null"`
}

// func (Employee) TableName() string {
// 	return "employees"
// }

// HashPassword converts plain text passwords into bcrypt hashes
func (e *Employee) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	e.Password = string(bytes)
	return nil
}

// CheckPassword compares hashed password with plain text password
func (e *Employee) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(password))
}

// BeforeCreate is a GORM hook that is executed before insert
func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	// if password not null dan not been hash, do hash
	if e.Password != "" && len(e.Password) < 60 {
		return e.HashPassword(e.Password)
	}
	return nil
}