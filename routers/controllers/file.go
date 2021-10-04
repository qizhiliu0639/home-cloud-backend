package controllers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"home-cloud/service"
	"net/http"
	"strings"
)

func UploadFiles(c *gin.Context) {
	//Todo Use dynamic routing to pass customer path for uploading
	form, _ := c.MultipartForm()
	fmt.Println(form)

	//field name for uploading form
	files := form.File["file"]
	session := sessions.Default(c)
	username, ok := session.Get("user").(string)
	if !ok || len(username) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": 1})
		return
	}
	//Todo Exit gracefully when uploading duplicate files
	for _, file := range files {
		if err := service.UploadFile(file, username, c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": 1})
			return
		}

	}
	c.JSON(http.StatusOK, gin.H{"success": 0})
}

func GetFileOrFolder(c *gin.Context) {
	path := c.Param("path")
	paths := strings.Split(path, "/")
	session := sessions.Default(c)
	username, ok := session.Get("user").(string)
	if !ok || len(username) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": 1})
		return
	}
	err := service.GetFileOrFolder(paths, username, c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": 1, "message": err})
	}

}
