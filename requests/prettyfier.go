package requests

import (
	"fmt"
	"strconv"
	"strings"
)

//CanteenToString returns a human readable string for a single canteen instance
func CanteenToString(canteen *Canteen, multiLine bool) string {
	if multiLine {
		return fmt.Sprintf("ID: %d\nName: %s\nCity: %s\nAddress: %s", canteen.ID, canteen.Name, canteen.City, canteen.Address)
	}
	return fmt.Sprintf("ID: %d\tName: %s\tCity: %s\tAddress: %s", canteen.ID, canteen.Name, canteen.City, canteen.Address)
}

//CanteenListToString returns a human readable string for a list of canteens
func CanteenListToString(canteens []Canteen) string {
	builder := strings.Builder{}
	for _, canteen := range canteens {
		builder.WriteString(CanteenToString(&canteen, false) + "\n")
	}
	return builder.String()
}

func priceToString(price prices, showOnlyStudent bool, seperator string) string {
	builder := strings.Builder{}
	if showOnlyStudent == true {
		builder.WriteString("Price: ")
		builder.WriteString(fmt.Sprintf("students %0.2f€", price.Students))
		return builder.String()
	}

	builder.WriteString("Prices: ")
	builder.WriteString(fmt.Sprintf("students %0.2f€,%s", price.Students, seperator))
	builder.WriteString(fmt.Sprintf("pupils %0.2f€,%s", price.Pupils, seperator))
	builder.WriteString(fmt.Sprintf("employees %0.2f€,%s", price.Employees, seperator))
	builder.WriteString(fmt.Sprintf("others %0.2f€%s", price.Others, seperator))
	return builder.String()
}

//CanteenMealToString returns a human readable string for a single canteenmeal instance
func CanteenMealToString(meal *CanteenMeal, seperator string, showPrice bool, showCategory bool, showNotes bool, showOnlyStudent bool) string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("Meal: %s", meal.Name) + seperator)

	if showCategory {
		builder.WriteString(fmt.Sprintf("Categorie: %s", meal.Category) + seperator)
	}

	if showNotes {
		builder.WriteString(fmt.Sprintf("Notes: %s", meal.Notes) + seperator)
	}

	if showPrice {
		builder.WriteString(fmt.Sprintf("%s", priceToString(meal.Prices, showOnlyStudent, " ")))
	}

	return builder.String()
}

//CanteenMealListToString returns a human readable string for a list if canteenmeals
func CanteenMealListToString(meals []CanteenMeal, showPrice, showNotes, showCategory, showOnlyStudent bool) string {
	builder := strings.Builder{}
	for i, meal := range meals {
		builder.WriteString(strconv.Itoa(i+1) + " " + CanteenMealToString(&meal, "\t", showPrice, showCategory, showNotes, showOnlyStudent) + "\n")
	}
	return builder.String()
}

//CanteenMealWeekListToString returns a human readable string for a whole week of meals
func CanteenMealWeekListToString(cateenWeek []CanteenDate, mealweek [][]CanteenMeal, showPrice, showNotes, showCategory, showOnlyStudent bool) string {
	builder := strings.Builder{}
	for i := range mealweek {
		if i == 0 {
			builder.WriteString(fmt.Sprintf("-> %s\n", cateenWeek[i].Date))
		} else {
			builder.WriteString(fmt.Sprintf("\n-> %s\n", cateenWeek[i].Date))
		}

		for _, meal := range mealweek[i] {
			builder.WriteString(fmt.Sprintf("\t%s\n", CanteenMealToString(&meal, "\t", showPrice, showCategory, showNotes, showOnlyStudent)))
		}
	}
	return builder.String()
}
