package main

import (
	"flag"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

func main() {

	// Set up db
	db, _ := gorm.Open("postgres", os.Getenv("DATABASE"))
	db.DB()
	err := db.DB().Ping()
	if err != nil {
		log.Panic(err)
	}
	db.SingularTable(true)

	// Continue with designated action and pass along db
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
