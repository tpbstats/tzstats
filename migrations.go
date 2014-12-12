package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func migrations(db gorm.DB) {

	db.LogMode(true)

	fmt.Println("Are you sure?")
	var answer string
	fmt.Scanf("%s", &answer)
	if answer != "yes" {
		return
	}

	db.Exec("drop schema public cascade")
	db.Exec("create schema public")
	db.AutoMigrate(&Scrape{}, &Status{})
}
