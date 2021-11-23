package utils

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Config struct {
	UserDataPath string `json:"user_data_path"`
	DBHost       string `json:"db_host"`
	DBPort       string `json:"db_port"`
	DBUser       string `json:"db_user"`
	DBPassword   string `json:"db_password"`
	DBName       string `json:"db_name"`
}

var globalConfig *Config
var configOnce sync.Once

func loadConfig() {
	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("Read config.json error. Please try to remove it or check the permission. ")
	}
	err = json.Unmarshal(jsonFile, &globalConfig)
	if err != nil {
		panic("Parse config.json error: " + err.Error() +
			" Please remove it and run again to generate a new one. ")
	}
	if globalConfig.UserDataPath == "" ||
		globalConfig.DBHost == "" ||
		globalConfig.DBPort == "" ||
		globalConfig.DBUser == "" ||
		globalConfig.DBPassword == "" ||
		globalConfig.DBName == "" {
		panic("Parse config.json error: missing some required fields. " +
			"Please check the file or generate a new one by removing it and running again. ")
	}
}

func GetConfig() *Config {
	configOnce.Do(loadConfig)
	return globalConfig
}
