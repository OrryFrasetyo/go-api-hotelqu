package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := os.Getenv("MYSQL_URL")

	log.Printf("DEBUG: MYSQL_URL value from environment: %s", dsn)
	// Pastikan dsn tidak kosong
	if dsn == "" {
		log.Println("ERROR: MYSQL_URL environment variable is empty!")
		panic("MYSQL_URL environment variable is not set. Please check Railway variables.")
	}

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("ERROR: Failed to connect to database using DSN: %s", dsn)
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