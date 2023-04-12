package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Finance *gorm.DB

func init() {
	Init("/Users/tuocheng/go/src/github.com/tuoaitang/calculator/finance/finance.db")
}

func Init(dbName string) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	Finance = db
}
