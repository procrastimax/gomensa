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
		fmt.Println(canteens)
		fmt.Println("canteen count:", len(canteens))
	}
}
