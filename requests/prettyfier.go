package requests

import (
	"fmt"
	"strconv"
	"strings"
)

//CanteenToString returns a human readable string for a single canteen instance
func CanteenToString(canteen *Canteen) string {
	return fmt.Sprintf("ID: %d\n\tName: %s\n\tCity: %s\n\tAddress: %s\n", canteen.ID, canteen.Name, canteen.City, canteen.Address)
}

//CanteenListToString returns a human readable string for a list of canteens
func CanteenListToString(canteens []Canteen) string {
	builder := strings.Builder{}
	for _, canteen := range canteens {
		builder.WriteString(CanteenToString(&canteen) + "\n")
	}
	return builder.String()
}

func notesToString(notes []string) string {
	builder := strings.Builder{}

	for _, note := range notes {
		builder.WriteString(fmt.Sprintf("\t\t- %s\n", note))
	}
	return builder.String()
}

func priceToString(price prices, showOnlyStudent bool, showOnlyEmployees bool, showOnlyOthers bool, showOnlyPupils bool, seperator string) string {
	builder := strings.Builder{}

	if showOnlyStudent == true {
		builder.WriteString("\n\tPrice:\n")
		builder.WriteString(fmt.Sprintf("\t\t- students: %0.2f€", price.Students))
		return builder.String()

	} else if showOnlyEmployees == true {
		builder.WriteString("\n\tPrice:\n")
		builder.WriteString(fmt.Sprintf("\t\t- employees: %0.2f€", price.Employees))
		return builder.String()

	} else if showOnlyOthers == true {
		builder.WriteString("\n\tPrice:\n")
		builder.WriteString(fmt.Sprintf("\t\t- others: %0.2f€", price.Others))
		return builder.String()

	} else if showOnlyPupils == true {
		builder.WriteString("\n\tPrice:\n")
		builder.WriteString(fmt.Sprintf("\t\t- pupils: %0.2f€", price.Pupils))
		return builder.String()
	}

	builder.WriteString("\n\tPrices:\n")
	builder.WriteString(fmt.Sprintf("\t\t- students: %0.2f€\n", price.Students))

	//only show pupils value, when its not 0
	if price.Pupils != 0.0 {
		builder.WriteString(fmt.Sprintf("\t\t- pupils: %0.2f€\n", price.Pupils))
	}

	builder.WriteString(fmt.Sprintf("\t\t- employees: %0.2f€\n", price.Employees))
	builder.WriteString(fmt.Sprintf("\t\t- others: %0.2f€\n", price.Others))
	return builder.String()
}

//CanteenMealToString returns a human readable string for a single canteenmeal instance
func CanteenMealToString(meal *CanteenMeal, showPrice bool, showCategory bool, showNotes bool, showOnlyStudent bool, showOnlyEmployees bool, showOnlyOthers bool, showOnlyPupils bool) string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("Meal: %s", meal.Name))

	if showCategory {
		if !showNotes && !showPrice {
			builder.WriteString(fmt.Sprintf("\n\tCategorie:\n\t\t- %s\n", meal.Category))
		} else {
			builder.WriteString(fmt.Sprintf("\n\tCategorie:\n\t\t- %s", meal.Category))
		}
	}

	if showNotes {
		builder.WriteString(fmt.Sprintf("\n\tNotes:\n%s", notesToString(meal.Notes)))
	}

	if showPrice {
		builder.WriteString(fmt.Sprintf("%s", priceToString(meal.Prices, showOnlyStudent, showOnlyEmployees, showOnlyOthers, showOnlyPupils, " ")))
	}

	builder.WriteString("\n")
	return builder.String()
}

//CanteenMealListToString returns a human readable string for a list if canteenmeals
func CanteenMealListToString(canteenDate CanteenDate, meals []CanteenMeal, canteen *Canteen, showPrice, showNotes, showCategory, showOnlyStudent bool, showOnlyEmployees bool, showOnlyOthers bool, showOnlyPupils bool) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s meals for date: %s:\n", canteen.Name, canteenDate.Date))
	for i, meal := range meals {
		builder.WriteString(strconv.Itoa(i+1) + " " + CanteenMealToString(&meal, showPrice, showCategory, showNotes, showOnlyStudent, showOnlyEmployees, showOnlyOthers, showOnlyPupils))
	}
	return builder.String()
}

//CanteenMealWeekListToString returns a human readable string for a whole week of meals
func CanteenMealWeekListToString(canteenWeek []CanteenDate, mealweek [][]CanteenMeal, canteen *Canteen, showPrice, showNotes, showCategory, showOnlyStudent, showOnlyEmployees, showOnlyOthers, showOnlyPupils bool) string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("%s meals for dates: %s - %s\n", canteen.Name, canteenWeek[0].Date, canteenWeek[len(canteenWeek)-1].Date))

	for i := range mealweek {
		if i == 0 {
			builder.WriteString("-> " + canteenWeek[i].Date + ":\n")
		} else {
			builder.WriteString("\n-> " + canteenWeek[i].Date + ":\n")
		}

		for j, meal := range mealweek[i] {
			builder.WriteString(fmt.Sprintf("%d %s", j+1, CanteenMealToString(&meal, showPrice, showCategory, showNotes, showOnlyStudent, showOnlyEmployees, showOnlyOthers, showOnlyPupils)))
		}
	}
	return builder.String()
}

//CanteenDateOpenedToString returns a string date representation of a canteen date with its opening information
func CanteenDateOpenedToString(canteenDate *CanteenDate, canteenName string, showWeek bool) string {
	builder := strings.Builder{}

	if showWeek == false && len(canteenName) > 1 {
		builder.WriteString(canteenName)
		builder.WriteString(" is open or closed on the following date:\n")
	}

	builder.WriteString(" - " + canteenDate.Date)
	if canteenDate.Closed {
		builder.WriteString(" -> closed")
	} else {
		builder.WriteString(" -> open")
	}
	return builder.String()
}

//CanteenDateListToString returns a prettified version of a list of canteen dates
func CanteenDateListToString(canteenDates []CanteenDate, canteenName string) string {
	builder := strings.Builder{}

	builder.WriteString(canteenName)
	builder.WriteString(" is open or closed on the following dates:\n")

	for _, date := range canteenDates {
		builder.WriteString(fmt.Sprintf("\t%s\n", CanteenDateOpenedToString(&date, "", false)))
	}
	return builder.String()
}
