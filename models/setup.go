package models

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("MYSQL_PUBLIC_URL")

	// Pastikan dsn tidak kosong
	if dsn == "" {
		panic("MYSQL_URL environment variable is not set")
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database : %v", err))
	}

	fmt.Println("Starting database migration...")

	err = database.AutoMigrate(&Department{}, &Position{}, &Shift{}, &Employee{}, &Schedule{}, &Attendance{})
	if err != nil {
			panic("failed to migrate Department, Position, Shift, Employee, " + err.Error())
	}
	
	fmt.Println("Database migration completed successfully")

	DB = database
}