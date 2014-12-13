package main

import (
	"fmt"
)

func migrations() {

	fmt.Println("Are you sure?")
	var answer string
	fmt.Scanf("%s", &answer)
	if answer != "yes" {
		return
	}

	db.LogMode(true)
	db.Exec("drop schema public cascade")
	db.Exec("create schema public")
	db.AutoMigrate(&Scrape{}, &Status{}, &Torrent{}, &Movie{}, &Category{})
}
