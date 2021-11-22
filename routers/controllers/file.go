package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"home-cloud/models"
	"home-cloud/service"
	"home-cloud/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

func UploadFiles(c *gin.Context) {
	form, _ := c.MultipartForm()

	user := c.Value("user").(*models.User)

	files := form.File["file"]
	vDir := c.Value("vDir").([]string)

	folder, err := service.GetFileOrFolderInfoByPath(vDir, user)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
		return
	}

	res := make(map[string]interface{})
	for _, file := range files {
		if len(file.Filename) == 0 || strings.ContainsAny(file.Filename, "/?*|<>:\\") {
			res[file.Filename] = gin.H{
				"result":  false,
				"message": "Invalid File Name",
			}
		} else {
			if err := service.UploadFile(file, user, folder, c); err != nil {
				res[file.Filename] = gin.H{
					"result":  false,
					"message": GetErrorMessage(err),
				}
			} else {
				res[file.Filename] = gin.H{
					"result": true,
				}
			}
		}
	}
	// If entering uploading process, success will be always 0 and each file result will be in files array
	c.JSON(http.StatusOK, gin.H{
		"success": 0,
		"files":   res,
	})
}

func GetFolder(c *gin.Context) {
	user := c.Value("user").(*models.User)
	vDir := c.Value("vDir").([]string)

	folder, err := service.GetFileOrFolderInfoByPath(vDir, user)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
		return
	}
	var files []*models.File

	files, err = service.GetFolder(folder, user)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0, "children": files})
	}
}

func NewFileOrFolder(c *gin.Context) {
	user := c.Value("user").(*models.User)
	vDir := c.Value("vDir").([]string)

	folder, err := service.GetFileOrFolderInfoByPath(vDir, user)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
		return
	}
	newName := c.PostForm("name")
	t := c.PostForm("type")
	if !(t == "file" || t == "folder") {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Type"})
		return
	}
	if len(newName) == 0 || strings.ContainsAny(newName, "/?*|<>:\\") {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Name"})
		return
	}
	err = service.NewFileOrFolder(folder, user, newName, t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}

func GetFile(c *gin.Context) {
	//This will only return error page in plain text because it may not be processed by axios
	user := c.Value("user").(*models.User)
	vDir := c.Value("vDir").([]string)

	file, err := service.GetFileOrFolderInfoByPath(vDir, user)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOrPermission) {
			c.String(http.StatusNotFound, "404 Not Found")
		} else if errors.Is(err, service.ErrSystem) {
			c.String(http.StatusInternalServerError, "500 Internal Server Error")
		} else {
			c.String(http.StatusBadRequest, "400 Bad Request")
		}
		return
	}
	var dst string
	var filename string
	dst, filename, err = service.GetFile(file, user)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOrPermission) {
			c.String(http.StatusNotFound, "404 Not Found")
		} else {
			c.String(http.StatusBadRequest, "400 Bad Request")
		}
		return
	}
	var f []byte
	f, err = ioutil.ReadFile(dst)
	if err != nil {
		utils.GetLogger().Errorf("Error when finding %s for %s", dst, file.Position)
		c.String(http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Content-Length", fmt.Sprintf("%d", len(f)))
	c.Header("Content-Type", http.DetectContentType(f))
	_, err = c.Writer.Write(f)
	if err != nil {
		//Delete header for download
		c.Writer.Header().Del("Content-Disposition")
		c.Writer.Header().Del("Content-Length")
		c.Writer.Header().Del("Content-Type")
		utils.GetLogger().Errorf("Error when writing %s to response", dst)
		c.String(http.StatusInternalServerError, "500 Internal Server Error")
	}
}

func GetFileOrFolderInfoByPath(c *gin.Context) {
	user := c.Value("user").(*models.User)
	vDir := c.Value("vDir").([]string)
	file, err := service.GetFileOrFolderInfoByPath(vDir, user)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		if file.IsDir == 1 {
			c.JSON(http.StatusOK, gin.H{"success": 0, "type": "folder", "root": file.ParentId == uuid.Nil, "info": file})
		} else {
			var folder *models.File
			folder, err = service.GetFileOrFolderInfoByID(file.ParentId, user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": 1, "message": GetErrorMessage(err)})
			} else {
				c.JSON(http.StatusOK, gin.H{"success": 0, "type": "file", "info": file, "parent_root": file.ParentId == uuid.Nil, "parent_info": folder})
			}
		}
	}

}

func DeleteFile(c *gin.Context) {
	user := c.Value("user").(*models.User)
	vDir := c.Value("vDir").([]string)

	file, err := service.GetFileOrFolderInfoByPath(vDir, user)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
		return
	}

	err = service.DeleteFile(file, user)
	//Will not raise error after starting to delete files
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}

func DealWithFavorite(c *gin.Context) {
	user := c.Value("user").(*models.User)
	vDir := c.Value("vDir").([]string)

	file, err := service.GetFileOrFolderInfoByPath(vDir, user)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
		return
	}

	err = service.ChangeFavoriteStatus(file, user)
	if err != nil {
		var status int
		if errors.Is(err, service.ErrInvalidOrPermission) {
			status = http.StatusNotFound
		} else if errors.Is(err, service.ErrSystem) {
			status = http.StatusInternalServerError
		} else {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}

}
