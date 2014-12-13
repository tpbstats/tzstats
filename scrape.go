package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/lib/pq"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var base string = "https://torrentz.eu/"
var rex map[string]*regexp.Regexp = map[string]*regexp.Regexp{
	"imdb": regexp.MustCompile(`imdb.com/title/(tt\d+)`),
	"cats": regexp.MustCompile(`».*$`),
}

func scrape() {

	db.LogMode(true)

	log.Println("Scrape commencing")

	urls := make([]string, 1)
	for i := 0; i < len(urls); i++ {
		urls[i] = fmt.Sprintf("%sany?q=movies&p=%d", base, i)
	}

	var documents = make([]*goquery.Document, len(urls))
	for key, url := range urls {
		document, err := getDocument(url)
		if err != nil {
			log.Panicln(err)
		}
		documents[key] = document
	}

	log.Println("Scrape")
	scrape := Scrape{}
	db.Save(&scrape)

	log.Println("Statuses")
	set := make(map[string]bool)
	for _, document := range documents {
		lists := document.Find(".results dl:not(:last-of-type)")
		lists.EachWithBreak(func(i int, list *goquery.Selection) bool {

			href, _ := list.Find("dt a").Attr("href")
			hash := href[1:]
			if _, exists := set[hash]; exists {
				return true
			}
			set[hash] = true

			term := list.Find("dt")
			cats := strings.Fields(rex["cats"].FindString(term.Text())[3:])
			torrent := Torrent{Hash: hash}
			if db.Find(&torrent, torrent).RecordNotFound() {
				torrent = scrapeTorrent(hash, cats)
				db.Save(&torrent)
			}

			status := Status{
				Seeders:   stringToInt(list.Find("dd span.u").Text()),
				Leechers:  stringToInt(list.Find("dd span.d").Text()),
				TorrentId: torrent.Id,
				ScrapeId:  scrape.Id,
			}
			db.Save(&status)

			return false
		})
	}

	log.Printf("Scrape complete, id=%d", scrape.Id)
}

func scrapeTorrent(hash string, cats []string) Torrent {

	torrent := Torrent{Hash: hash}

	for _, cat := range cats {
		torrent.Categories = append(torrent.Categories, Category{Name: cat})
	}

	url := fmt.Sprintf("%s%s", base, torrent.Hash)
	document, err := getDocument(url)
	if err != nil {
		log.Println(err)
		return torrent
	}

	links := document.Find(".download dl:not(:first-of-type) dt a")
	urls := make([]string, 0, links.Length())
	links.Each(func(i int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if !exists {
			return
		}
		urls = append(urls, href)
	})
	torrent.Movie = scrapeMovie(urls)

	torrent.Rating, _ = strconv.Atoi(document.Find(".votebox .status").Text())

	return torrent
}

func scrapeMovie(urls []string) Movie {
	movie := Movie{}
	for _, url := range urls {
		body, err := getBody(url)
		if err != nil {
			continue
		}
		matches := rex["imdb"].FindStringSubmatch(body)
		if matches == nil || len(matches) < 2 {
			continue
		}
		movie.Imdb = matches[1]
		break
	}
	return movie
}
