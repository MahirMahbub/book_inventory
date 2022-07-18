package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	dbHost := os.Getenv("DB_HOST_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbTimeZone := os.Getenv("DB_TIMEZONE")
	dbSSLMode := os.Getenv("DB_SSL")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimeZone)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	err = database.AutoMigrate(&User{})
	if err != nil {
		return
	}
	err = database.AutoMigrate(&Token{})
	if err != nil {
		return
	}
	err = database.AutoMigrate(&NonAuthToken{})
	if err != nil {
		return
	}
	DB = database
}
