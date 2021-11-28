package bootstrap

import (
	"encoding/json"
	"fmt"
	"home-cloud/models"
	"home-cloud/utils"
	"io/ioutil"
	"os"
	"strings"
)

// BootStrap check config file and perform init operations before running
// use in bootstrap package to prevent import cycle
func BootStrap() {
	configFile, err := os.Stat("config.json")
	if err != nil {
		// If it is a file not exist error, create the template file
		if os.IsNotExist(err) {
			initConfigJson()
			fmt.Println(
				"Find config.json error. " +
					"We have generated an example with following default settings. " +
					"You can choose NO to modify these settings.",
			)
			fmt.Println()
			fmt.Println("User Data Path: ./data")
			fmt.Println("Mysql Host: 127.0.0.1")
			fmt.Println("Mysql Port: 3306")
			fmt.Println("Mysql Username: root")
			fmt.Println("Mysql Password: root")
			fmt.Println("Mysql Database Name: homecloud")
			fmt.Println("Listen Address: 127.0.0.1:8080")
			var input string
			for {
				fmt.Print("Do you want to continue? [Y/n] ")
				_, err = fmt.Scanln(&input)
				if err != nil {
					panic("User input error: " + err.Error())
				} else {
					input = strings.ToLower(input)
					if input == "y" {
						break
					} else if input == "n" {
						os.Exit(0)
					}
				}
			}
		} else {
			panic("Read config file error: " + err.Error())
		}
	} else {
		if configFile.IsDir() {
			panic("Please remove the config.json directory before running")
		}
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
	cfg.ListenAddress = "127.0.0.1:8080"
	jsonFile, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		panic("Create config.json error: " + err.Error())
	}
	err = ioutil.WriteFile("config.json", jsonFile, 0644)
	if err != nil {
		panic("Create config.json error: " + err.Error())
	}
}
