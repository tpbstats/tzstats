package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var client http.Client = http.Client{
	Timeout: time.Duration(2 * time.Second),
}

func getBody(url string, attempts int) (string, error) {
	for i := 0; i < attempts; i++ {
		resp, err := client.Get(url)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		return string(body), nil
	}
	return "", fmt.Errorf("Failure %s", url)
}

func getDocument(url string, attempts int) (*goquery.Document, error) {
	for i := 0; i < attempts; i++ {
		resp, err := client.Get(url)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		document, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			continue
		}
		return document, nil
	}
	return nil, fmt.Errorf("Failure %s", url)
}

func stringToInt(str string) int {
	str = strings.Replace(str, ",", "", -1)
	number, _ := strconv.Atoi(str)
	return number
}
