package main

import (
	"home-cloud/models"
	"home-cloud/routers"
	"home-cloud/utils"
)

func main() {
	models.InitDatabase()
	router := routers.InitRouter()
	//router.Use(middleware.LoggerToFile())
	err := router.Run()
	if err != nil {
		utils.GetLogger().Panic("Error to run!")
	}
}
