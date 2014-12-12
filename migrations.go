package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func migrations(db gorm.DB) {

	db.LogMode(true)

	// Check if sure
	fmt.Println("Are you sure?")
	var answer string
	fmt.Scanf("%s", &answer)
	if answer != "yes" { return }

	// Define tables
	tables := []interface{}{
		&Status{},
		&Scrape{},
	}

	// Drop tables
	for i := len(tables)-1; i >= 0; i-- {
		db.DropTableIfExists(tables[i])
	}

	// Create tables
	for _, table := range tables {
		db.CreateTable(table)
	}
}
