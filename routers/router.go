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
		{
			dirGroup := fileAPI.Group("")
			dirGroup.Use(middleware.ValidateDir())
			{
				//Upload file
				dirGroup.POST("/upload", controllers.UploadFiles)
				//Get child in folder (Use folder ID)
				dirGroup.POST("/list_dir", controllers.GetFolder)
				//New file or Folder
				dirGroup.POST("/new", controllers.NewFileOrFolder)

				dirGroup.POST("/get_info", controllers.GetFileOrFolderInfoByPath)
				//Get file (Use file name)
				dirGroup.POST("/get_file", controllers.GetFile)
				//delete file
				dirGroup.POST("/delete", controllers.DeleteFile)
				//Add favorite file
				dirGroup.PUT("/favorite", controllers.ToggleFavorite)
			}
			//Search file by keywords
			fileAPI.POST("/search", controllers.SearchFiles)
			//Get Favorites List
			fileAPI.GET("/get_favorite", controllers.GetFavorites)
		}
		userAPI := api.Group("/user")
		userAPI.Use(middleware.AuthSession())
		{
			userAPI.PUT("/password", controllers.ChangePassword)
			userAPI.POST("/profile", controllers.UpdateProfile)
		}
	}
}
