package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"home-cloud/service"
	"home-cloud/utils"
	"net/http"
)

// UserPreLogin Frontend will request this before login to get the salt for PBKDF2
func UserPreLogin(c *gin.Context) {
	username := c.PostForm("username")
	if len(username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid username"})
	}
	accountSalt := service.LoginGetSalt(username)
	c.JSON(http.StatusOK, gin.H{"success": 0, "account_salt": accountSalt})
}

// UserLogin Login controller
func UserLogin(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid username or password"})
		return
	}

	if service.LoginValidate(username, password) {
		session := sessions.Default(c)
		session.Set("user", username)
		if err := session.Save(); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": 1, "message": "Save session error"})
			return
		}
		utils.GetLogger().Info("User " + username + " successfully log in")
		c.JSON(http.StatusOK, gin.H{"success": 0})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"success": 1, "message": "Username or password error"})
	}
}

// UserLogout Logout user
func UserLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	if err := session.Save(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Save session error"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}

// UserRegister Register User
func UserRegister(c *gin.Context) {
	//获取表单信息
	username := c.PostForm("username")
	password := c.PostForm("password")
	accountSalt := c.PostForm("accountSalt")

	if err := service.RegisterUser(username, password, accountSalt); err != nil {
		//Todo return specific information if it is a system error
		c.JSON(http.StatusForbidden, gin.H{"success": 1, "message": "Username invalid!"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}
