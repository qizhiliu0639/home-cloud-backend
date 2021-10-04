package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"home-cloud/routers/controllers"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("bacsadhkwqidh23@#$@#CA*Y(qada213411qwfe23!@$!R!@CASasdh1212CQAWF"),
		[]byte("akckq3213QWE!@EW!ESVGFQrqfaw23QW")) //cookie secret
	router.Use(sessions.Sessions("home-cloud-backend-session", store))

	api := router.Group("/api")
	{
		//Register
		api.POST("/register", controllers.UserRegister)
		//Login and logout
		//pre-Login for getting salt of the user, will return a random salt if user not exists
		api.POST("/pre-login", controllers.UserPreLogin)
		api.POST("/login", controllers.UserLogin)
		api.GET("/logout", controllers.UserLogout)

		//Files API
		fileAPI := api.Group("/file")
		{
			//Upload Files
			fileAPI.POST("/*path", controllers.UploadFiles)
			fileAPI.GET("/*path", controllers.GetFileOrFolder)
		}
	}
	return router
}
