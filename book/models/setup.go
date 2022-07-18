package models

import (
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	//dsn := "postgres://root:root@10.5.5.5:5432/test"
	//dsn := "host=10.5.5.5 user=root password=root dbname=test port=5432 sslmode=prefer"
	dsn := "host=10.5.5.5 user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	err = database.AutoMigrate(&Book{}, &Author{})
	if err != nil {
		return
	}
	DB = database
}
