package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"home-cloud/middleware"
	"home-cloud/routers/controllers"
)

func InitRouter(router *gin.Engine) {
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

		statusAPI := api.Group("/status")
		statusAPI.Use(middleware.AuthSession())
		{
			statusAPI.GET("/user", controllers.GetUserStatus)
		}

		//Files API
		fileAPI := api.Group("/file")
		fileAPI.Use(middleware.AuthSession())
		fileAPI.Use(middleware.ValidateDir())
		{
			//Upload file
			fileAPI.POST("/upload", controllers.UploadFiles)
			//Get child in folder (Use folder ID)
			fileAPI.POST("/list_dir", controllers.GetFolder)
			//New file or Folder
			fileAPI.POST("/new", controllers.NewFileOrFolder)

			fileAPI.POST("/get_info", controllers.GetFileOrFolderInfoByPath)
			//Get file (Use file ID)
			fileAPI.POST("/get_file", controllers.GetFile)
			//delete file
			fileAPI.POST("/delete", controllers.DeleteFile)
			//Add favorite file
			fileAPI.PUT("/favorite", controllers.DealWithFavorite)
		}
	}
}
