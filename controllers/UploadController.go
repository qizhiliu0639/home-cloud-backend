package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"home-cloud/models"
	"log"
	"net/http"
	"time"
)

const dst = "./Files/"

func UploadGet(c *gin.Context){

	files,_ := models.ListAllFiles()

	c.HTML(http.StatusOK, "upload.html", gin.H{"title": "upload page","File":files})
}

func UploadSingleFile(c *gin.Context){
	file, _ := c.FormFile("upload")
	log.Println(file.Filename)

	// 上传文件至指定目录
	c.SaveUploadedFile(file, dst+file.Filename)

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))

	filemodel := models.File{0, dst+"/"+file.Filename, file.Filename, 0, time.Now().Unix()}
	models.InsertFile(filemodel)
	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}


func UploadMultiFiles(c *gin.Context){
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		log.Println(file.Filename)

		// 上传文件至指定目录
		c.SaveUploadedFile(file, dst+file.Filename)
		filemodel := models.File{0, dst+file.Filename, file.Filename, 0, time.Now().Unix()}
		models.InsertFile(filemodel)
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

func responseErr(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": err})
}