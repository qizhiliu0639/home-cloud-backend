package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"home-cloud/bootstrap"
	"home-cloud/middleware"
	"home-cloud/routers"
	"io/ioutil"
	"log"
)

//go:embed web/build
var frontendFS embed.FS

func main() {
	bootstrap.BootStrap()
	// disable log in gin session middleware (it uses the standard logger)
	log.SetOutput(ioutil.Discard)
	gin.SetMode(gin.ReleaseMode)
	// disable log color
	gin.DisableConsoleColor()
	router := gin.Default()
	router.Use(middleware.FrontendFileHandler(frontendFS, "web/build"))
	routers.InitRouter(router)
	err := router.Run("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
}
