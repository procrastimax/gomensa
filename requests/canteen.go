package requests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	openMensaEndpoint = "https://openmensa.org/api/v2"
)

// abort handles the abortion when requesting all available canteens
var abort = make(chan struct{})

// sema shall limit the number of goroutines for requesting all available canteens
var sema = make(chan struct{}, 5)

//Canteen is a struct representing a single canteen instance without geopgrapical coordinates
type Canteen struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Address string `json:"address"`
}

//RequestCanteenByID makes a get request for retrieving a single Canteen by its ID
func RequestCanteenByID(ID uint32) *Canteen {
	baseURL, err := url.Parse(openMensaEndpoint)
	if err != nil {
		log.Println("ERROR: Malformed URL ", err.Error())
		return nil
	}

	// Add a Path Segment (Path segment is automatically escaped)
	baseURL.Path += "/canteens/" + strconv.Itoa(int(ID))

	resp, err := http.Get(baseURL.String())
	if err != nil {
		log.Println("ERROR: Something went wrong when requesting a list of all canteens!", err.Error())
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR: Something went wrong when processing the result of requesting a list of all canteens!", err.Error())
		return nil
	}

	var canteen Canteen
	err = json.Unmarshal(body, &canteen)
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to parse the requestcanteen result!", err.Error())
		return nil
	}

	return &canteen
}

//RequestListOfAllCanteens request all canteens from all api pages and return a list of all
func RequestListOfAllCanteens() []Canteen {
	//currently there are more than 400 canteens, so we can allocate some memory before appending the slices
	allCanteens := make([]Canteen, 0, 400)

	canteensChan := make(chan []Canteen)
	page := 1

loop:
	for {
		select {
		case sema <- struct{}{}: //acquire token
		case <-abort:
			break loop
		}

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
	defer func() { <-sema }()
	baseURL, err := url.Parse(openMensaEndpoint)
	if err != nil {
		log.Println("ERROR: Malformed URL ", err.Error())
		abort <- struct{}{}
		return
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
		abort <- struct{}{}
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR: Something went wrong when processing the result of requesting a list of all canteens!", err.Error())
		abort <- struct{}{}
		return
	}

	var canteens []Canteen
	err = json.Unmarshal(body, &canteens)
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to parse the canteen list results!", err.Error())
		abort <- struct{}{}
		return
	}

	//cleaning random new lines
	for i := range canteens {
		canteens[i].Name = strings.ReplaceAll(canteens[i].Name, "\n", "")
		canteens[i].City = strings.ReplaceAll(canteens[i].City, "\n", "")
		canteens[i].Address = strings.ReplaceAll(canteens[i].Address, "\n", "")
	}

	canteensChan <- canteens

	// check if the next page would be the last page and then closes the channel
	maxPages, err := strconv.Atoi(resp.Header.Get("X-Total-Pages"))
	if err != nil {
		log.Println("ERROR: Could not convert max page header key to int!")
		abort <- struct{}{}
		return
	}

	if page+1 > maxPages {
		abort <- struct{}{}
		return
	}
}
