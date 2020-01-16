package configutil

import (
	"encoding/json"
	"gomensa/requests"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

const (
	//configFilePath is the path of the config file
	configFilePath = "/.config/gomensa/"
	//configFileName is the name of the config file itself
	configFileName = "config.json"

	//rwPermissionPath permission bits for reading, writing, executing -> used for config gomensa folder
	rwPermissionPath = 0700
	//rwPermissionFile permission bits for reading, writing of config file
	rwPermissionFile = 0644
)

//Config represents the user settings of which canteen he usually visits for eating
type Config struct {
	Canteen requests.Canteen `json:"canteen"`
}

//SaveConfig saves a user configuration to the .config/gomensa folder, the config file is saved as a json file
// returns a bool value indicating the success of the file save process
func SaveConfig(config *Config) bool {
	//if config file doesnt exist -> create it
	if CheckConfigExists() == false {
		//create folder if config file does not exist
		err := os.MkdirAll(getUsersHomeDir()+configFilePath, rwPermissionPath)
		if err != nil {
			log.Println("ERROR: Something went wrong when trying to create the necessary folders for config file!", err.Error())
			return false
		}
	}
	json, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to convert Config struct to json object!", err.Error())
		return false
	}

	err = ioutil.WriteFile(getUsersHomeDir()+configFilePath+configFileName, json, rwPermissionFile)
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to save json object to file!", err.Error())
		return false
	}

	return true
}

//ReadConfig reads the config file from the .config/gomensa/ directory, returns nil when the directory/ config file does not exist
func ReadConfig() *Config {
	if CheckConfigExists() == false {
		log.Println("ERROR: Can't read config file because file does not exist!")
		log.Println("Using empty config file.")
		return nil
	}

	configContent, err := ioutil.ReadFile(getUsersHomeDir() + configFilePath + configFileName)
	if err != nil {
		log.Println("ERROR: Something went wrong trying to open existing config file!", err.Error())
		return nil
	}

	config := &Config{}
	err = json.Unmarshal(configContent, config)
	if err != nil {
		log.Println("ERROR: Something went wrong when trying to parse config file to json object!", err.Error())
		return nil
	}
	return config
}

//CheckConfigExists checks whether the config file exists and returns a bool
func CheckConfigExists() bool {
	_, err := os.Stat(getUsersHomeDir() + configFilePath + configFileName)

	if os.IsNotExist(err) {
		return false
	}

	if err == nil {
		return true
	}

	log.Println("ERROR: Something went wrong when checking the existance of the config file!", err.Error())
	return false
}

//getUsersHomeDir returns the users home directory as string like /home/USER
func getUsersHomeDir() string {
	myself, err := user.Current()
	if err != nil {
		log.Println("ERROR: Can't get users home directory!")
	}
	homedir := myself.HomeDir
	return homedir
}
