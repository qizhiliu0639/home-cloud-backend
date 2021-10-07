package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"home-cloud/middleware"
	"home-cloud/models"
	"home-cloud/routers"
)

//go:embed web/build
var frontendFS embed.FS

func main() {
	models.InitDatabase()
	router := gin.Default()
	router.Use(middleware.FrontendFileHandler(frontendFS, "web/build"))
	routers.InitRouter(router)
	err := router.Run("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
}
