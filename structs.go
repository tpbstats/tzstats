package main

import (
	"time"
)

type Scrape struct {
	Id        int64
	CreatedAt time.Time
}

type Status struct {
	Id        int64
	Seeders   int
	Leechers  int
	Torrent   Torrent
	TorrentId int64
	ScrapeId  int64
}

type Torrent struct {
	Id         int64
	Hash       string `sql:"type:char(40);unique"`
	Rating     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Movie      Movie
	MovieId    int64
	Categories []Category `gorm:"many2many:torrent_category;"`
}

type Movie struct {
	Id        int64
	Imdb      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Category struct {
	Id   int64
	Name string
}
