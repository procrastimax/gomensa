package tests

import (
	"fmt"
	"gomensa/requests"
	"testing"
)

func TestRequestAllCanteens(t *testing.T) {
	canteens := requests.RequestListOfAllCanteens()
	if len(canteens) == 0 {
		t.Error("Did not retrieve any canteens when requesting a list of them!")
	} else {
		fmt.Println("canteen count:", len(canteens))
	}
}

func TestRequestCanteenByID(t *testing.T) {
	canteen := requests.RequestCanteenByID(1)
	if canteen == nil {
		t.Error("Could not retrieve single canteen by ID!")
	} else {
		fmt.Println(*canteen)
	}

	id := 6666
	canteen = requests.RequestCanteenByID(uint32(id))
	if canteen != nil {
		t.Errorf("This canteen with ID %d should not exist!", id)
	}
}

func TestVariousCanteenDates(t *testing.T) {
	canteenWeek := requests.RequestCanteenWeek(32)
	if len(canteenWeek) == 0 {
		t.Error("Something went wrong, we sould definetly retrieve a week of canteen dates here!")
	} else {
		fmt.Println(canteenWeek)
	}

	canteenDay := requests.RequestCanteenDateToday(32)
	if len(canteenDay.Date) == 0 {
		t.Error("Something went wrong, we sould definetly retrieve a canteenDate for today!")
	} else {
		fmt.Println(*canteenDay)
	}

	canteenDay = requests.RequestCanteenDateTomorrow(32)
	if len(canteenDay.Date) == 0 {
		t.Error("Something went wrong, we sould definetly retrieve a canteenDate for tomorrow!")
	} else {
		fmt.Println(*canteenDay)
	}
}

func TestVariousCanteenMealDates(t *testing.T) {
	_, canteenMeals := requests.RequestCanteenMealOfToday(32)
	if len(canteenMeals) == 0 {
		t.Error("Could not retrieve the meals for today!")
	}

	_, canteenMeals = requests.RequestCanteenMealOfTomorrow(32)
	if len(canteenMeals) == 0 {
		t.Error("Could not retrieve the meals for tomorrow!")
	}

	_, canteenWeekMeals := requests.RequestCanteenMealsOfWeek(32)
	if len(canteenWeekMeals) == 0 {
		t.Error("Could not retrieve the meals for week!")
	}
}
