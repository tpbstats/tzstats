package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func scrape(db gorm.DB) {

	// Config
	pages := 1
	attempts := 3

	// Get urls
	urls := make([]string, pages)
	base := "https://torrentz.eu/search?f=movie&p="
	for i := 0; i < len(urls); i++ {
		urls[i] = fmt.Sprintf("%s%d", base, i)
	}

	// Get responses
	var responses = make([]*http.Response, len(urls))
	client := http.Client{
	    Timeout: time.Duration(2 * time.Second),
	}
	for key, url := range urls {
		for i := 0; i < attempts; i++ {
			resp, err := client.Get(url)
			if err != nil || resp.StatusCode != 200 {
				log.Printf("error: %s", url)
				continue
			}
			responses[key] = resp
			break
		}
	}

	// Get statuses
	scrape := Scrape{
		Time:     time.Now(),
		Statuses: make([]Status, len(responses)*50),
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
			key := i*50 + j
			status := Status{
				Hash:     href[1:],
				Seeders:  seeders,
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
