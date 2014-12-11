package main

import (
	"flag"
	"github.com/jinzhu/gorm"
	"os"
	"log"
)

func main() {
	var action string
	flag.StringVar(&action, "action", "scrape", "action")
	flag.Parse()

	db, err := gorm.Open("postgres", os.Getenv("DATABASE"))
	if err != nil {
		log.Panic(err)
	}
	db.DB()
	db.SingularTable(true)
	db.LogMode(true)

	switch action {
	case "migrations":
		migrations(db)
	case "scrape":
		scrape(db)
	default:
		panic("Invalid action")
	}
}
