package tests

import (
	"fmt"
	"gomensa/configutil"
	"gomensa/requests"
	"testing"
)

func TestSaveConfig(t *testing.T) {
	config := configutil.Config{
		Canteen: requests.Canteen{
			ID:      0,
			Name:    "test",
			City:    "test",
			Address: "test",
		},
	}
	configutil.SaveConfig(&config)
	fmt.Println("Saved ", config)
	if configutil.CheckConfigExists() == false {
		t.Error("SaveConfig failed when trying to create config file!")
	}
}

func TestReadConfig(t *testing.T) {
	config := configutil.ReadConfig()
	if config == nil {
		t.Error("Could not read config from files!")
	}
}
