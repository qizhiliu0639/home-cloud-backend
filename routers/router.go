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

		//Files API
		fileAPI := api.Group("/file")
		fileAPI.Use(middleware.AuthSession())
		{
			dir := fileAPI.Group("")
			dir.Use(middleware.ValidateID("dir"))
			{
				//Upload file
				dir.POST("/upload", controllers.UploadFiles)
				//Get child in folder (Use folder ID)
				dir.POST("/list_dir", controllers.GetFolder)
				//New file or Folder
				dir.POST("/new", controllers.NewFileOrFolder)
			}
			fileID := fileAPI.Group("")
			fileID.Use(middleware.ValidateID("fileID"))
			{
				fileID.POST("/get_info_id", controllers.GetFileOrFolderInfoByID)
				//Get file (Use file ID)
				fileID.POST("/get_file", controllers.GetFile)
				//delete file
				fileID.POST("/delete", controllers.DeleteFile)
			}
			//Check if path exists and return metadata.
			fileAPI.POST("/get_info_path", controllers.GetFileOrFolderInfoByPath)
		}
	}
}
