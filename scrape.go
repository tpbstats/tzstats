package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"time"
)

func scrape(db gorm.DB) {

	// Base urls
	bases := []string{
		"https://torrentz.eu/",
		"https://torrentz.me/",
		"https://torrentz.ch/",
		"https://torrentz.in/",
	}

	// Relative page urls
	relatives := []string{
		"search?f=movie&p=0",
		// "search?f=movie&p=1",
		// "search?f=movie&p=2",
	}

	// Get responses
	var responses = make([]*http.Response, len(relatives))
	for key, relative := range relatives {
		for _, base := range bases {
			url := fmt.Sprintf("%s%s", base, relative)
			resp, err := http.Get(url)
			if err != nil {
				log.Println("Problem connecting")
				continue
			}
			responses[key] = resp
			break
		}
	}

	// Get statuses
	scrape := Scrape{
		Time: time.Now(),
		Statuses: make([]Status, len(responses) * 50),
	}
	for i, response := range responses {
		document, err := goquery.NewDocumentFromResponse(response)
		if err != nil {
			log.Println("Failed to create document")
			continue
		}
		lists := document.Find(".results dl:not(:last-of-type)")
		lists.Each(func(j int, list *goquery.Selection) {
			href, _ := list.Find("dt a").Attr("href")
			seeders := stringToInt(list.Find("dd span.u").Text())
			leechers := stringToInt(list.Find("dd span.d").Text())
			key := i * 50 + j
			status := Status{
				Hash: href[1:],
				Seeders: seeders,
				Leechers: leechers,
			}
			scrape.Statuses[key] = status
		})
	}

	// Iterate through statusses to get torrents

	// Insert scrape
	db.Create(&scrape)
}

func stringToInt(str string) int {
	str = strings.Replace(str, ",", "", -1)
	number, _ := strconv.Atoi(str)
	return number
}