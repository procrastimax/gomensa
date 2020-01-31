package requests

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
	dateMatchingString = "\\d\\d\\d\\d-\\d\\d-\\d\\d"

	pageFlag      = 1
	limitFlag     = 2
	startDateFlag = 4
)

var (
	dateMatchingRegex = regexp.MustCompile(dateMatchingString)
)

//CanteenDate is a struct representing a date of a single canteen and if the canteen is closed at this date
type CanteenDate struct {
	Date   string `json:"date"`
	Closed bool   `json:"closed"`
}

//RequestCanteenDateTomorrow calls the requestDatesOfCanteen function with the limit = 1, a page = 2 and no startDate, so we retrieve the canteen date of tomorrow
func RequestCanteenDateTomorrow(ID uint32) (*CanteenDate, bool) {
	canteenDay := requestDatesOfCanteen(ID, "", 2, 1)
	if canteenDay == nil || len(canteenDay) == 0 {
		return &CanteenDate{}, false
	}
	return &canteenDay[0], true
}

//RequestCanteenDateToday calls the requestDatesOfCanteen function with the limit of 1 and no startDate, so we retrieve the current date as a canteen date
func RequestCanteenDateToday(ID uint32) (*CanteenDate, bool) {
	canteenDay := requestDatesOfCanteen(ID, "", 0, 1)
	if canteenDay == nil || len(canteenDay) == 0 {
		return &CanteenDate{}, false
	}
	return &canteenDay[0], true
}

//RequestCanteenDate given the ID and a date in the format YYYY-MM-DD the an instance of canteendate is returned with information about whether or not the canteen is opened on this day
func RequestCanteenDate(ID uint32, date string) (*CanteenDate, bool) {
	canteenDate := requestDatesOfCanteen(ID, date, 0, 1)

	if canteenDate == nil || len(canteenDate) == 0 {
		return &CanteenDate{}, false
	}

	return &canteenDate[0], true
}

//RequestCanteenWeek calls the requestDatesOfCanteen function with the limit of 7 and no startDate, so we retrieve the next 7 days of a canteen
func RequestCanteenWeek(ID uint32) ([]CanteenDate, bool) {
	canteenWeek := requestDatesOfCanteen(ID, "", 0, 7)
	//when a nil slice was passed, it is most likely an error occured, so we return a false ok value
	if canteenWeek == nil || len(canteenWeek) == 0 {
		return []CanteenDate{}, false
	}
	return canteenWeek, true
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

	var canteenDates []CanteenDate
	err = json.Unmarshal(body, &canteenDates)
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to parse the canteen list results!", err.Error())
		return nil
	}
	return canteenDates
}
