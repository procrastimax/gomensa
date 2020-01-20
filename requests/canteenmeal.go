package requests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

//CanteenMeal is a struct representing a single meal of a canteen
type CanteenMeal struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Notes    []string `json:"notes"`
	Prices   prices   `json:"prices"`
	Category string   `json:"category"`
}

type prices struct {
	Students  float64 `json:"students"`
	Employees float64 `json:"employees"`
	Pupils    float64 `json:"pupils"`
	Others    float64 `json:"others"`
}

//RequestCanteenMealOfTomorrow returns all meals that are offered at the given canteen tomorrow
func RequestCanteenMealOfTomorrow(canteenID uint32) (*CanteenDate, []CanteenMeal) {
	canteenDateToday, ok := RequestCanteenDateTomorrow(canteenID)

	if canteenDateToday.Date == "" || ok == false {
		log.Println("Something went wrong when trying to request mensa date of tomorrow!")
		return &CanteenDate{}, []CanteenMeal{}
	}

	canteenMeals := requestCanteenMeals(canteenID, canteenDateToday.Date)
	if canteenMeals == nil {
		return &CanteenDate{}, []CanteenMeal{}
	}
	return canteenDateToday, canteenMeals
}

//RequestCanteenMealsOfWeek returns all meals of the next 7 days from a given canteen
func RequestCanteenMealsOfWeek(canteenID uint32) ([]CanteenDate, [][]CanteenMeal) {
	canteenDateList, ok := RequestCanteenWeek(canteenID)

	if len(canteenDateList) == 0 || ok == false {
		log.Println("Something went wrong when trying to request mensa dates of week!")
		return []CanteenDate{}, [][]CanteenMeal{}
	}

	canteenMealList := make([][]CanteenMeal, len(canteenDateList))
	mealList := []CanteenMeal{}

	for i, date := range canteenDateList {
		mealList = requestCanteenMeals(canteenID, date.Date)

		if mealList == nil {
			canteenMealList[i] = nil
			break
		}
		canteenMealList[i] = mealList
	}
	return canteenDateList, canteenMealList
}

//RequestCanteenMealOfToday returns the canteenMeal for the current day
//this functions makes a requestCanteenDate request to see if the canteen is open and if there is any information provided about the meals
func RequestCanteenMealOfToday(canteenID uint32) (*CanteenDate, []CanteenMeal) {
	canteenDateToday, ok := RequestCanteenDateToday(canteenID)
	//check if we retrieved an empty instance of the canteenDate
	if canteenDateToday.Date == "" || ok == false {
		log.Println("Something went wrong when trying to request mensa date of today!")
		return &CanteenDate{}, []CanteenMeal{}
	}

	canteenMeals := requestCanteenMeals(canteenID, canteenDateToday.Date)
	if canteenMeals == nil {
		log.Println("Something went wrong when trying to request mensa meal of today!")
		return &CanteenDate{}, []CanteenMeal{}
	}
	return canteenDateToday, canteenMeals
}

//requestCanteenMeals is a function which requests a list of meals for a given canteen with a canteenID and a canteenDate
func requestCanteenMeals(canteendID uint32, canteenDate string) []CanteenMeal {
	baseURL, err := url.Parse(openMensaEndpoint)
	if err != nil {
		log.Println("ERROR: Malformed URL ", err.Error())
		return nil
	}

	// Add a Path Segment (Path segment is automatically escaped)
	baseURL.Path += "/canteens/" + strconv.Itoa(int(canteendID)) + "/days/" + canteenDate + "/meals"

	resp, err := http.Get(baseURL.String())
	if err != nil {
		log.Println("ERROR: Something went wrong when requesting a list of meals!", err.Error())
		return nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ERROR: Something went wrong when processing the result of requesting a list of all meals for a day!", err.Error())
		return nil
	}

	var canteenMeals []CanteenMeal
	err = json.Unmarshal(body, &canteenMeals)
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to parse the requestCantenenMeal result!", err.Error())
		return nil
	}

	return canteenMeals
}
