package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func scrape(db gorm.DB) {

	// Base urls
	bases := []string{
		"http://torrentz.eu/",
		"http://torrentz.me/",
		"http://torrentz.ch/",
		"http://torrentz.in/",
	}

	// Relative page urls
	relatives := make([]string, 10)
	for i := 0; i < len(relatives); i++ {
		relatives[i] = fmt.Sprintf("search?f=movie&p=%d", i)
	}

	// Get responses
	transport := http.Transport{
		Dial: func(network, url string) (net.Conn, error) {
			return net.DialTimeout(network, url, 2*time.Second)
		},
	}
	client := http.Client{
		Transport: &transport,
	}
	var responses = make([]*http.Response, len(relatives))
	for key, relative := range relatives {
		for _, base := range bases {
			url := fmt.Sprintf("%s%s", base, relative)
			resp, err := client.Get(url)
			if err != nil || resp.StatusCode != 200 {
				log.Printf("error: %s", url)
				continue
			}
			responses[key] = resp
			log.Printf("success: %s", url)
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
