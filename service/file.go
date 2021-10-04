package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"home-cloud/models"
	"home-cloud/utils"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
)

func generateRealPath(filename string) string {
	//Todo Generate UUID and recover if conflicts or panics
	return filename
}

func UploadFile(upFile *multipart.FileHeader, username string, c *gin.Context) error {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return errors.New("user not exists")
	}
	file := models.NewFile()
	file.IsDir = 0
	file.Name = upFile.Filename
	file.OwnerId = user.ID
	file.CreatorId = user.ID
	file.RealPath = generateRealPath(upFile.Filename)
	file.Size = uint64(upFile.Size)
	//Todo Upload to different folder based on parameter
	var rootFolder *models.File
	rootFolder, err = user.GetRootFolder()
	if err != nil {
		return errors.New("get folder error")
	}
	file.ParentId = rootFolder.ID
	err = file.CreateFile()
	if err != nil {
		return errors.New("create file record error")
	}

	dst := path.Join(utils.GetConfig().UserDataPath, strconv.FormatUint(user.ID, 10),
		"data", "files", file.RealPath)
	utils.GetLogger().Infof("Save file to %s", dst)
	if err = c.SaveUploadedFile(upFile, dst); err != nil {
		return errors.New("save upFile error")
	}
	return nil
}

func GetFileOrFolder(path []string, username string, c *gin.Context) error {
	user, err := models.GetUserByUsername(username)
	if err != nil {
		return errors.New("user not exists")
	}
	var rootFolder *models.File
	rootFolder, err = user.GetRootFolder()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	//Only for root path now
	if len(path) == 2 && path[0] == ""&&path[1]=="" {
		files, err := rootFolder.GetChildInFolder()
		if err != nil {
			fmt.Println("Here")

			return err
		}
		c.JSON(http.StatusOK, files)
	}
	return nil
}
