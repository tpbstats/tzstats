package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func scrape(db gorm.DB) {

	log.Println("Commencing")

	// Get urls
	urls := make([]string, 1)
	base := "https://torrentz.eu/search?f=movie&p="
	for i := 0; i < len(urls); i++ {
		urls[i] = fmt.Sprintf("%s%d", base, i)
	}

	// Get responses
	var responses = make([]*http.Response, len(urls))
	for key, url := range urls {
		log.Printf("Getting %s", url)
		resp, err := get(url)
		if err != nil {
			log.Panicln(err)
			continue
		}
		responses[key] = resp
		log.Printf("Success %s", url)
	}

	// Get statuses
	log.Println("Statuses")
	scrape := Scrape{
		Time:     time.Now(),
		Statuses: make([]Status, len(responses)*50),
	}
	for i, response := range responses {
		document, err := goquery.NewDocumentFromResponse(response)
		if err != nil {
			log.Panicln(err)
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
	log.Println("Inserting")
	db.Create(&scrape)
	log.Printf("Inserted, id=%d", scrape.Id)

	log.Println("Finished")
}
