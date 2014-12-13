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

	urls := make([]string, 10)
	base := "https://torrentz.eu/search?f=movie&p="
	for i := 0; i < len(urls); i++ {
		urls[i] = fmt.Sprintf("%s%d", base, i)
	}

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

	log.Println("Statuses")
	statuses := make(map[string]Status)
	for _, response := range responses {
		document, err := goquery.NewDocumentFromResponse(response)
		if err != nil {
			log.Panicln(err)
			continue
		}
		lists := document.Find(".results dl:not(:last-of-type)")
		lists.Each(func(i int, list *goquery.Selection) {
			href, _ := list.Find("dt a").Attr("href")
			status := Status{}
			status.Hash = href[1:]
			status.Seeders = stringToInt(list.Find("dd span.u").Text())
			status.Leechers = stringToInt(list.Find("dd span.d").Text())
			statuses[status.Hash] = status
		})
	}

	log.Println("Assembling")
	scrape := Scrape{
		Time:     time.Now(),
		Statuses: make([]Status, 0, len(statuses)),
	}
	for _, status := range statuses {
		scrape.Statuses = append(scrape.Statuses, status)
	}

	log.Println("Inserting")
	db.Create(&scrape)
	log.Printf("Inserted, id=%d", scrape.Id)

	log.Println("Finished")
}
