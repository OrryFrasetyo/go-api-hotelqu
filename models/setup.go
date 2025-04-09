package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:@tcp(127.0.0.1:3306)/hotelqu_db?parseTime=true"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Starting database migration...")

	err = database.AutoMigrate(&Department{}, &Position{}, &Shift{}, &Employee{}, &Schedule{}, &Attendance{})
	if err != nil {
			panic("failed to migrate Department, Position, Shift, Employee, " + err.Error())
	}
	
	fmt.Println("Database migration completed successfully")

	DB = database
}