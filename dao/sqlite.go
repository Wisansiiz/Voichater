package dao

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DBlite *gorm.DB

func InitSqlite() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DBlite = db
}
