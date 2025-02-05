package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DBConnection *gorm.DB

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DBConnection = db

	return DBConnection
}

func GetDB() *gorm.DB {
	return DBConnection
}
