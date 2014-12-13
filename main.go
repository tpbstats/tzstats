package main

import (
	"flag"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var db gorm.DB

func main() {

	log.Println("Connecting to database")
	db, _ = gorm.Open("postgres", os.Getenv("DATABASE"))
	db.DB()
	defer db.DB().Close()
	err := db.DB().Ping()
	if err != nil {
		log.Panic(err)
	}
	db.SingularTable(true)
	log.Println("Connection established")

	// Continue with designated action
	var action string
	flag.StringVar(&action, "action", "scrape", "action")
	flag.Parse()
	switch action {
	case "migrations":
		migrations()
	case "scrape":
		scrape()
	default:
		panic("Invalid action")
	}
}
