package routers

import(
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"home-cloud/controllers"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")
	store := cookie.NewStore([]byte("loginuser"))
	router.Use(sessions.Sessions("mysession", store))
	{
		//注册：
		router.GET("/register",controllers.RegisterGet)
		router.POST("/register",controllers.RegisterPost)

		//登录
		router.GET("/login",controllers.LoginGet)
		router.POST("/login",controllers.LoginPost)

		router.GET("/exit", controllers.ExitGet)

		router.GET("/upload",controllers.UploadGet)
		router.POST("/uploadSingle",controllers.UploadSingleFile)
		router.POST("/uploadMulti",controllers.UploadMultiFiles)
	}
	return router

}
