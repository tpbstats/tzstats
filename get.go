package main

import (
	"fmt"
	"net/http"
	"time"
)

func get(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Duration(2 * time.Second),
	}
	for i := 0; i < 3; i++ {
		resp, err := client.Get(url)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Failed to get %s", url)
}