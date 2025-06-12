package models

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbUser := os.Getenv("MYSQL_USERNAME")  
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")     
	dbPort := os.Getenv("MYSQL_PORT")     
	dbName := os.Getenv("MYSQL_DATABASE")  
	
	// dsn := "root:@tcp(127.0.0.1:3306)/hotelqu_db?parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

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