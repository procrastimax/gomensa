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
	allCanteens := make([]Canteen, 100)

	page := 1
	for page <= 10 {
		requestCanteens(page)
		page++
	}

	return allCanteens
}

//requestCanteens makes a GET request to the openmensa endpoint and returns a list of all canteens
func requestCanteens(page int) []Canteen {
	if page <= 0 {
		log.Println("ERROR: Cannot use page parameter less than 1!")
		return []Canteen{}
	}

	baseURL, err := url.Parse(openMensaEndpoint)
	if err != nil {
		log.Println("ERROR: Malformed URL ", err.Error())
		return []Canteen{}
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
		return []Canteen{}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("ERROR: Something went wrong when processing the result of requesting a list of all canteens!", err.Error())
		return []Canteen{}
	}

	canteens := make([]Canteen, 0)
	err = json.Unmarshal(body, &canteens)

	if err != nil {
		log.Println("ERROR: Something went wrong when trying to parse the canteen list results!", err.Error())
		return []Canteen{}
	}

	return canteens
}
