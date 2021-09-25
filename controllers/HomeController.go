package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeGet(c *gin.Context) {
	//获取session，判断用户是否登录
	islogin := GetSession(c)
	c.HTML(http.StatusOK, "home.html", gin.H{"IsLogin": islogin})
}
