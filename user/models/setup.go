package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=10.5.5.5 user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
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
