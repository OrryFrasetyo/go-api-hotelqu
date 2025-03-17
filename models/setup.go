package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:@tcp(127.0.0.1:3306)/hotelqu_db"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Cetak pesan saat mulai migrasi
	fmt.Println("Starting database migration...")

	// database.AutoMigrate(&Department{}, &Position{}, &Employee{})

	// Migrasi model Department dan Position terlebih dahulu
	err = database.AutoMigrate(&Department{}, &Position{})
	if err != nil {
			panic("failed to migrate Department and Position: " + err.Error())
	}
	
	// Migrasi model Employee
	err = database.AutoMigrate(&Employee{})
	if err != nil {
			panic("failed to migrate Employee: " + err.Error())
	}
	
	fmt.Println("Database migration completed successfully")

	DB = database
}