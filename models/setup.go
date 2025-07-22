package models

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// setup for local
	dsn := "root:@tcp(127.0.0.1:3306)/hotelqu_db?parseTime=true"
	// rawURL := os.Getenv("mysql_url")
	// log.Printf("DEBUG: MYSQL_URL value from environment: %s", rawURL)

	// if rawURL == "" {
	// 	log.Println("ERROR: MYSQL_URL environment variable is empty!")
	// 	panic("MYSQL_URL is not set. Please set it in Railway.")
	// }

	// // Parse MYSQL_URL from railway
	// parsedURL, err := url.Parse(rawURL)
	// if err != nil {
	// 	panic(fmt.Sprintf("Invalid MYSQL_URL format: %v", err))
	// }

	// user := parsedURL.User.Username()
	// pass, _ := parsedURL.User.Password()
	// host := parsedURL.Hostname()
	// port := parsedURL.Port()
	// if port == "" {
	// 	port = "3306"
	// }
	// dbName := strings.TrimPrefix(parsedURL.Path, "/")

	// // Create DSN in Go-compatible format
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbName)

	// log.Printf("DEBUG: Final DSN used to connect: %s", dsn)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("ERROR: Failed to connect to database using DSN: %s", dsn)
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	fmt.Println("Starting database migration...")
	err = database.AutoMigrate(&Department{}, &Position{}, &Shift{}, &Employee{}, &Schedule{}, &Attendance{})
	if err != nil {
		panic("failed to migrate: " + err.Error())
	}

	fmt.Println("Database migration completed successfully")
	DB = database
}
