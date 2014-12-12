package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var client http.Client = http.Client{
	Timeout: time.Duration(2 * time.Second),
}

func get(url string) (*http.Response, error) {
	for i := 0; i < 3; i++ {
		resp, err := client.Get(url)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get %s", url)
}

func stringToInt(str string) int {
	str = strings.Replace(str, ",", "", -1)
	number, _ := strconv.Atoi(str)
	return number
}
