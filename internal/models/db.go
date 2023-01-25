package models

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("DB_DEV_DSN")

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}

func Migrate() {
	DB.AutoMigrate(&Group{}, &Member{}, &Schedule{})
}
