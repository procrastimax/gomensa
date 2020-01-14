package main

import (
	"testing"
)

func TestRequestAllCanteens(t *testing.T) {
	canteens := requestCanteens(1)

	if len(canteens) == 0 {
		t.Error("Did not retrieve any canteens when requesting a list of them!")
	}
}
