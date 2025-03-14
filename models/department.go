package models

type Department struct {
	Id                 int        `json:"id" gorm:"primary_key"`
	ParentDepartmentId *int       `json:"parent_department_id" gorm:"index"`
	DepartmentName     string     `json:"department_name"`
}