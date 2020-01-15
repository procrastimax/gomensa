package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const openMensaEndpoint = "https://openmensa.org/api/v2"

//Canteen is a struct representing a single canteen instance without geopgrapical coordinates
type Canteen struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Address string `json:"address"`
}

//RequestListOfAllCanteens request all canteens from all api pages and return a list of all
func RequestListOfAllCanteens() []Canteen {
	//currently there are more than 496 canteens, so we can allocate some memory before appending the slices -> Im not gonna update this value in the future
	allCanteens := make([]Canteen, 0, 496)

	canteensChan := make(chan []Canteen)
	page := 1
	for {
		go requestCanteens(page, canteensChan)
		page++

		value, ok := <-canteensChan
		//request and parse responses until the channel is closed -> this happens when we reach the max page number
		if ok == false {
			break
		}
		allCanteens = append(allCanteens, value...)
	}
	return allCanteens
}

//requestCanteens makes a GET request to the openmensa endpoint and returns a list of all canteens
func requestCanteens(page int, canteensChan chan<- []Canteen) {
	baseURL, err := url.Parse(openMensaEndpoint)
	if err != nil {
		log.Println("ERROR: Malformed URL ", err.Error())
		close(canteensChan)
	}

	// Add a Path Segment (Path segment is automatically escaped)
	baseURL.Path += "/canteens"

	// Prepare Query Parameters
	params := url.Values{}
	params.Add("page", strconv.Itoa(page))

	// Add Query Parameters to the URL
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		log.Println("ERROR: Something went wrong when requesting a list of all canteens!", err.Error())
		close(canteensChan)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR: Something went wrong when processing the result of requesting a list of all canteens!", err.Error())
		close(canteensChan)
	}

	canteens := make([]Canteen, 0)
	err = json.Unmarshal(body, &canteens)
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to parse the canteen list results!", err.Error())
		close(canteensChan)
	}

	canteensChan <- canteens

	// check if the next page would be the last page and then closes the channel
	maxPages, err := strconv.Atoi(resp.Header.Get("X-Total-Pages"))
	if err != nil {
		log.Println("ERROR: Could not convert max page header key to int!")
		close(canteensChan)
	}

	if page+1 > maxPages {
		close(canteensChan)
	}
}
