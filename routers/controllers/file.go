package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"home-cloud/models"
	"home-cloud/service"
	"io/ioutil"
	"net/http"
	"strings"
)

func UploadFiles(c *gin.Context) {
	form, _ := c.MultipartForm()

	user := c.Value("user").(*models.User)

	files := form.File["file"]
	vDir := c.Value("dir").(uuid.UUID)

	res := make(map[string]interface{})
	success := false
	for _, file := range files {
		if len(file.Filename) == 0 || strings.ContainsAny(file.Filename, "/?*|<>:\\") {
			res[file.Filename] = gin.H{
				"result":  false,
				"message": "Invalid File Name",
			}
		} else {
			if err := service.UploadFile(file, user, vDir, c); err != nil {
				res[file.Filename] = gin.H{
					"result":  false,
					"message": GetErrorMessage(err),
				}
			} else {
				success = true
				res[file.Filename] = gin.H{
					"result": true,
				}
			}
		}
	}
	if success {
		c.JSON(http.StatusOK, gin.H{
			"success": 0,
			"files":   res,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": 1,
			"files":   res,
		})
	}

}

func GetFolder(c *gin.Context) {
	vDir := c.Value("dir").(uuid.UUID)
	user := c.Value("user").(*models.User)

	files, err := service.GetFolder(vDir, user)
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
	vDir := c.Value("dir").(uuid.UUID)
	newName := c.PostForm("name")
	t := c.PostForm("type")
	user := c.Value("user").(*models.User)
	if !(t == "file" || t == "folder") {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Type"})
		return
	}
	if len(newName) == 0 || strings.ContainsAny(newName, "/?*|<>:\\") {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Name"})
		return
	}
	err := service.NewFileOrFolder(vDir, user, newName, t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": GetErrorMessage(err)})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}

func GetFile(c *gin.Context) {
	user := c.Value("user").(*models.User)
	vFileID := c.Value("fileID").(uuid.UUID)
	dst, filename, err := service.GetFile(vFileID, user)
	if err != nil {
		if errors.Is(err, service.ErrInvalidOrPermission) {
			c.JSON(http.StatusNotFound, gin.H{"success": 1, "message": GetErrorMessage(err)})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": GetErrorMessage(err)})
		}
		return
	}
	var f []byte
	f, err = ioutil.ReadFile(dst)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": 1, "message": "Server Error"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"success": 1, "message": "Server Error"})
	}
}

func GetFileOrFolderInfoByID(c *gin.Context) {
	user := c.Value("user").(*models.User)

	vFileID := c.Value("fileID").(uuid.UUID)
	file, err := service.GetFileOrFolderInfoByID(vFileID, user)
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
		c.JSON(http.StatusOK, gin.H{"success": 0, "file": file})
	}
}

func GetFileOrFolderInfoByPath(c *gin.Context) {
	user := c.Value("user").(*models.User)
	path := c.PostForm("path")
	if len(path) < 1 || !strings.HasPrefix(path, "/") || strings.HasSuffix(path, "/") {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Path"})
		return
	}
	//filter root slash
	path = path[1:]
	var paths []string
	//not a root path
	if len(path) > 0 {
		pathsTmp := strings.Split(path, "/")
		for _, p := range pathsTmp {
			if len(p) < 1 || p == "." || p == ".." {
				c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Path"})
				return
			} else {
				paths = append(paths, p)
			}
		}
	}
	file, err := service.GetFileOrFolderInfoByPath(paths, user)
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
		c.JSON(http.StatusOK, gin.H{"success": 0, "file": file})
	}
}

func DeleteFile(c *gin.Context) {
	user := c.Value("user").(*models.User)
	vFileID := c.Value("fileID").(uuid.UUID)
	err := service.DeleteFile(vFileID, user)
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
