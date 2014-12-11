package main

import (
	"flag"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

func main() {

	// Set up DB
	db, err := gorm.Open("postgres", os.Getenv("DATABASE"))
	if err != nil {
		log.Panic(err)
	}
	db.DB()
	db.SingularTable(true)
	// db.LogMode(true)

	// Continue with designated action and pass along DB
	var action string
	flag.StringVar(&action, "action", "scrape", "action")
	flag.Parse()
	switch action {
	case "migrations":
		migrations(db)
	case "scrape":
		scrape(db)
	default:
		panic("Invalid action")
	}
}
