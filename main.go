package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/shiena/ansicolor"
	"home-cloud/bootstrap"
	"home-cloud/middleware"
	"home-cloud/routers"
	"os"
)

//go:embed web/build
var frontendFS embed.FS

func main() {
	bootstrap.BootStrap()
	gin.SetMode(gin.ReleaseMode)
	gin.ForceConsoleColor()
	gin.DefaultWriter = ansicolor.NewAnsiColorWriter(os.Stdout)
	router := gin.Default()
	router.Use(middleware.FrontendFileHandler(frontendFS, "web/build"))
	routers.InitRouter(router)
	err := router.Run("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
}
