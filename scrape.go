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

	// Get urls
	urls := make([]string, 1)
	base := "https://torrentz.eu/search?f=movie&p="
	for i := 0; i < len(urls); i++ {
		urls[i] = fmt.Sprintf("%s%d", base, i)
	}

	// Get responses
	var responses = make([]*http.Response, len(urls))
	for key, url := range urls {
		resp, err := get(url)
		if err != nil || resp.StatusCode != 200 {
			log.Printf("error: %s", url)
			continue
		}
		responses[key] = resp
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

	// Insert scrape
	db.Create(&scrape)
}

func stringToInt(str string) int {
	str = strings.Replace(str, ",", "", -1)
	number, _ := strconv.Atoi(str)
	return number
}
