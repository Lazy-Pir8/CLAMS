package core

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the database
var DB *gorm.DB

// connecting to postgres
func ConnectDB() {
	dsn := "host=localhost user=postgres password=mysecretpassword dbname=postgres port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database Connection Sucessful")
}
