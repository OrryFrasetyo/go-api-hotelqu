package models

type Position struct {
	Id           int     `json:"id" gorm:"primary_key"`
	DepartmentId int        `json:"department_id" gorm:"index"`
	Department   Department `json:"department" gorm:"foreignKey:DepartmentId"`
	PositionName string     `json:"position_name"`
	IsCompleted  bool       `json:"is_completed"`
}