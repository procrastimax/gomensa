package main

import (
	"fmt"
	"testing"
)

func TestRequestAllCanteens(t *testing.T) {
	canteens := RequestListOfAllCanteens()
	if len(canteens) == 0 {
		t.Error("Did not retrieve any canteens when requesting a list of them!")
	} else {
		fmt.Println("canteen count:", len(canteens))
	}
}

func TestRequestCanteenByID(t *testing.T) {
	canteen := RequestCanteenByID(1)
	if canteen == nil {
		t.Error("Could not retrieve single canteen by ID!")
	} else {
		fmt.Println(*canteen)
	}

	id := 6666
	canteen = RequestCanteenByID(uint32(id))
	if canteen != nil {
		t.Errorf("This canteen with ID %d should not exist!", id)
	}
}

func TestRequestDatesOfCanteens(t *testing.T) {
	canteenDates := requestDatesOfCanteen(31, "", 0, 0)
	if canteenDates == nil {
		t.Error("NIL returned, could not get valid dates of canteens 1")
	} else if len(canteenDates) == 0 {
		fmt.Println("Got an empty list of canteen dates! 1")
	}

	canteenDates = requestDatesOfCanteen(31, "", 3, 0)
	if canteenDates == nil {
		t.Error("NIL returned, could not get valid dates of canteens 2")
	} else if len(canteenDates) == 0 {
		fmt.Println("Got an empty list of canteen dates! 2")
	}

	canteenDates = requestDatesOfCanteen(31, "", 0, 3)
	if canteenDates == nil {
		t.Error("NIL returned, could not get valid dates of canteens 3")
	} else if len(canteenDates) == 0 {
		fmt.Println("Got an empty list of canteen dates! 3")
	}

	canteenDates = requestDatesOfCanteen(31, "2020-02-13", 3, 3)
	if canteenDates == nil {
		t.Error("NIL returned, could not get valid dates of canteens 4")
	} else if len(canteenDates) == 0 {
		fmt.Println("Got an empty list of canteen dates! 4")
	}

	canteenDates = requestDatesOfCanteen(321, "2020-02-13", 3, 3)
	if canteenDates == nil {
		t.Error("NIL returned, could not get valid dates of canteens 5")
	} else if len(canteenDates) == 0 {
		fmt.Println("Got an empty list of canteen dates! 5")
	}
}

func TestVariousCanteenDates(t *testing.T) {
	canteenWeek := RequestCanteenWeek(32)
	if len(canteenWeek) == 0 {
		t.Error("Something went wrong, we sould definetly retrieve a week of canteen dates here!")
	} else {
		fmt.Println(canteenWeek)
	}

	canteenDay := RequestCanteenDateToday(32)
	if len(canteenDay.Date) == 0 {
		t.Error("Something went wrong, we sould definetly retrieve a canteenDate for today!")
	} else {
		fmt.Println(*canteenDay)
	}

	canteenDay = RequestCanteenDateTomorrow(32)
	if len(canteenDay.Date) == 0 {
		t.Error("Something went wrong, we sould definetly retrieve a canteenDate for tomorrow!")
	} else {
		fmt.Println(*canteenDay)
	}
}

func TestRequestCanteenMeals(t *testing.T) {
	canteenMeals := requestCanteenMeals(32, "2020-01-15")
	if canteenMeals == nil {
		t.Error("Could not retrieve list of meals!")
	}
}

func TestVariousCanteenMealDates(t *testing.T) {
	canteenMeals := RequestCanteenMealOfToday(32)
	if len(canteenMeals) == 0 {
		t.Error("Could not retrieve the meals for today!")
	}

	canteenMeals = RequestCanteenMealOfTomorrow(32)
	if len(canteenMeals) == 0 {
		t.Error("Could not retrieve the meals for tomorrow!")
	}

	canteenWeekMeals := RequestCanteenMealsOfWeek(32)
	if len(canteenWeekMeals) == 0 {
		t.Error("Could not retrieve the meals for week!")
	}
}
