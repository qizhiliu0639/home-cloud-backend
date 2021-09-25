package main

import (
	"home-cloud/database"
	"home-cloud/middleware"
	"home-cloud/routers"
)

func main() {
	database.InitMysql()
	router := routers.InitRouter()
	router.Use(middleware.LoggerToFile())
	router.Static("/static", "./static")
	router.Run()
}
