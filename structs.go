package main

import(
	"time"
)

type Scrape struct {
	Id int
	Time time.Time
	Statuses []Status
}

type Status struct {
	Id int
	Hash string
	Seeders int
	Leechers int
	ScrapeId int
}