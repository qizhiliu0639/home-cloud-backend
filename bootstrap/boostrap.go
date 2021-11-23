package bootstrap

import (
	"encoding/json"
	"home-cloud/models"
	"home-cloud/utils"
	"io/ioutil"
	"os"
)

// BootStrap check config file and perform init operations before running
// use in bootstrap package to prevent import cycle
func BootStrap() {
	configFile, err := os.Stat("config.json")
	if err != nil {
		initConfigJson()
		panic(
			"Find config.json error:" + err.Error() +
				" We have generated an example. " +
				"Please change the fields in config.json and run again",
		)
	}
	if configFile.IsDir() {
		panic("Please remove the config.json directory before running")
	}

	models.InitDatabase()
}

func initConfigJson() {
	var cfg utils.Config
	cfg.UserDataPath = "./data"
	cfg.DBHost = "127.0.0.1"
	cfg.DBPort = "3306"
	cfg.DBUser = "root"
	cfg.DBPassword = "root"
	cfg.DBName = "homecloud"
	jsonFile, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		panic("Create config.json error: " + err.Error())
	}
	err = ioutil.WriteFile("config.json", jsonFile, 0644)
	if err != nil {
		panic("Create config.json error: " + err.Error())
	}
}
