package main

import (
	"fmt"
	"testing"
)

func TestSaveConfig(t *testing.T) {
	config := Config{
		Canteen: Canteen{
			ID:      0,
			Name:    "test",
			City:    "test",
			Address: "test",
		},
	}
	SaveConfig(&config)
	fmt.Println("Saved ", config)
	if checkConfigExists() == false {
		t.Error("SaveConfig failed when trying to create config file!")
	}
}

func TestReadConfig(t *testing.T) {
	config := ReadConfig()
	if config == nil {
		t.Error("Could not read config from files!")
	}
}
