package utils

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type config struct {
	LogFilePath  string `json:"log_file_path"`
	UserDataPath string `json:"user_data_path"`
	DBHost       string `json:"db_host"`
	DBPort       string `json:"db_port"`
	DBUser       string `json:"db_user"`
	DBPassword   string `json:"db_password"`
	DBName       string `json:"db_name"`
}

var globalConfig *config
var configOnce sync.Once

func loadConfig() {
	jsonFile, err := ioutil.ReadFile("config.json")
	if err != nil {
		//Todo Add default Config
		return
	}
	_ = json.Unmarshal(jsonFile, &globalConfig)
}

func GetConfig() *config {
	configOnce.Do(loadConfig)
	return globalConfig
}
