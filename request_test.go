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
