package main

import (
	"github.com/jinzhu/gorm"
)

func migrations(db gorm.DB) {
	db.AutoMigrate(&Scrape{}, &Status{})
}
