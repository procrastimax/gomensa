package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

const (
	openMensaEndpoint  = "https://openmensa.org/api/v2"
	dateMatchingString = "\\d\\d\\d\\d-\\d\\d-\\d\\d"

	pageFlag      = 1
	limitFlag     = 2
	startDateFlag = 4
)

var (
	dateMatchingRegex = regexp.MustCompile(dateMatchingString)
)

//Canteen is a struct representing a single canteen instance without geopgrapical coordinates
type Canteen struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	City    string `json:"city"`
	Address string `json:"address"`
}

//CanteenDate is a struct representing a date of a single canteen and if the canteen is closed at this date
type CanteenDate struct {
	Date   string `json:"date"`
	Closed bool   `json:"closed"`
}

//RequestCanteenDateTomorrow calls the requestDatesOfCanteen function with the limit = 1, a page = 2 and no startDate, so we retrieve the canteen date of tomorrow
func RequestCanteenDateTomorrow(ID uint32) *CanteenDate {
	canteenDay := requestDatesOfCanteen(ID, "", 2, 1)
	if canteenDay == nil {
		return &CanteenDate{}
	}
	return &canteenDay[0]
}

//RequestCanteenDateToday calls the requestDatesOfCanteen function with the limit of 1 and no startDate, so we retrieve the current date as a canteen date
func RequestCanteenDateToday(ID uint32) *CanteenDate {
	canteenDay := requestDatesOfCanteen(ID, "", 0, 1)
	if canteenDay == nil {
		return &CanteenDate{}
	}
	return &canteenDay[0]
}

//RequestCanteenWeek calls the requestDatesOfCanteen function with the limit of 7 and no startDate, so we retrieve the next 7 days of a canteen
func RequestCanteenWeek(ID uint32) []CanteenDate {
	canteenWeek := requestDatesOfCanteen(ID, "", 0, 7)
	if canteenWeek == nil {
		return []CanteenDate{}
	}
	return canteenWeek
}

//requestDatesOfCanteen requests Dates of canteens returning a list of CanteenDate for representing open/ closed dates of the canteen
//it is advised to expect that the returned list of dates can be empty, this is the case when to date information is given
//this function needs an ID, a startDate in the form YYYY-MM-DD for specifiyng a startDate for requesting, when passing an unvalid format or empty string the current date is used
//also a page and a maximal limit of returns can be specified if you dont need them set them to 0
func requestDatesOfCanteen(ID uint32, startDate string, page uint32, limit uint32) []CanteenDate {

	//useFlags is flag representing which parameters to use for requesting the canteendate
	usageFlags := uint8(0)

	if page != 0 {
		usageFlags |= pageFlag
	}

	if limit != 0 {
		usageFlags |= limitFlag
	}

	if len(startDate) > 0 {
		//only set startDateFlag when the date is valid!
		if dateMatchingRegex.MatchString(startDate) == false {
			log.Println("ERROR: The passed startDate does not follow the following structure: YYYY-MM-DD! Using current date.")
		} else {
			usageFlags |= startDateFlag
		}
	}

	baseURL, err := url.Parse(openMensaEndpoint)
	if err != nil {
		log.Println("ERROR: Malformed URL ", err.Error())
		return nil
	}

	// Add a Path Segment (Path segment is automatically escaped)
	baseURL.Path += "/canteens/" + strconv.Itoa(int(ID)) + "/days"

	// Prepare Query Parameters
	params := url.Values{}

	//use paging parameter
	if usageFlags&pageFlag != 0 {
		params.Add("page", strconv.Itoa(int(page)))
	}

	//use limit parameter
	if usageFlags&limitFlag != 0 {
		params.Add("limit", strconv.Itoa(int(limit)))
	}

	//use startDate parameter
	if usageFlags&startDateFlag != 0 {
		params.Add("start", startDate)
	}

	// Add Query Parameters to the URL
	baseURL.RawQuery = params.Encode()

	resp, err := http.Get(baseURL.String())
	if err != nil {
		log.Println("ERROR: Something went wrong when requesting a list of canteenDates!", err.Error())
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR: Something went wrong when processing the result of requesting a list of canteenDates!", err.Error())
		return nil
	}

	canteenDates := []CanteenDate{}
	err = json.Unmarshal(body, &canteenDates)
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to parse the canteen list results!", err.Error())
		return nil
	}
	return canteenDates
}

//RequestCanteenByID makes a get request for retrieving a single Canteen by its ID
func RequestCanteenByID(ID uint32) *Canteen {
	baseURL, err := url.Parse(openMensaEndpoint)
	if err != nil {
		log.Println("ERROR: Malformed URL ", err.Error())
		return nil
	}

	// Add a Path Segment (Path segment is automatically escaped)
	baseURL.Path += "/canteens"
	baseURL.Path += "/" + strconv.Itoa(int(ID))

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

	canteen := &Canteen{}
	err = json.Unmarshal(body, canteen)
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to parse the requestcanteen result!", err.Error())
		return nil
	}

	return canteen
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

	canteens := []Canteen{}
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
