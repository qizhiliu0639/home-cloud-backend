package controllers

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"home-cloud/models"
	"home-cloud/service"
	"home-cloud/utils"
	"net/http"
	"strconv"
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
	remember := c.PostForm("remember")
	encryptionKey := c.PostForm("encryptionKey")

	if len(username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid username or password"})
		return
	}
	validation, user := service.LoginValidate(username, password)

	if validation {
		if user.Migration == 1 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": 1,
				"message": "Migration in progress! You are not allowed to log in! ",
			})
			return
		} else if user.Migration == 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": 1,
				"message": "Migration error occurred! Please contact the admin! ",
			})
		}
		session := sessions.Default(c)
		session.Set("user", username)
		session.Set("encryptionKey", encryptionKey)
		session.Set("remember", remember)
		maxAge := 86400 * 30
		if remember == "0" {
			// maxAge=0 will delete cookie after browser close
			maxAge = 0
		}
		session.Options(sessions.Options{
			Path:     "/api",
			HttpOnly: true,
			MaxAge:   maxAge,
			SameSite: http.SameSiteLaxMode,
		})
		if err := session.Save(); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"success": 1, "message": "Save session error! "})
			return
		}
		utils.GetLogger().Info("User " + username + " successfully log in")
		c.JSON(http.StatusOK, gin.H{"success": 0})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"success": 1, "message": "Authentication error! "})
	}
}

// UserLogout Logout user
func UserLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	session.Delete("EncryptionKey")
	session.Delete("remember")
	session.Options(sessions.Options{
		Path:     "/api",
		HttpOnly: true,
		// MaxAge=0 will delete cookie after browser close
		MaxAge:   0,
		SameSite: http.SameSiteLaxMode,
	})
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
	encryptionKey := c.PostForm("encryption")

	if len(accountSalt) < 64 || len(password) < 64 {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Request!"})
		return
	}
	if err := service.RegisterUser(username, password, accountSalt, encryptionKey); err != nil {
		if errors.Is(err, service.ErrUsernameInvalid) {
			c.JSON(http.StatusForbidden, gin.H{"success": 1, "message": "Invalid Username!"})
		} else if errors.Is(err, service.ErrRequestPara) {
			c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": GetErrorMessage(err)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"success": 1, "message": "Server Error!"})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}

// ChangePassword Change user password
func ChangePassword(c *gin.Context) {
	user := c.Value("user").(*models.User)
	oldPassword := c.PostForm("old")
	newPassword := c.PostForm("new")
	newAccountSalt := c.PostForm("new_account_salt")
	newEncryption := c.PostForm("new_encryption")
	oldEncryption := c.PostForm("old_encryption")
	validation, _ := service.LoginValidate(user.Username, oldPassword)

	if validation {
		err := service.ChangePassword(user, newAccountSalt, newPassword, oldEncryption, newEncryption)
		if err != nil {
			if errors.Is(err, service.ErrRequestPara) {
				c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": GetErrorMessage(err)})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"success": 1, "message": "Server Error"})
			}
			return
		}
		session := sessions.Default(c)
		remember := session.Get("remember").(string)
		maxAge := 86400 * 30
		if remember == "0" {
			// maxAge=0 will delete cookie after browser close
			maxAge = 0
		}
		session.Set("user", user.Username)
		session.Set("encryptionKey", newEncryption)
		session.Set("remember", remember)
		session.Options(sessions.Options{
			Path:     "/api",
			HttpOnly: true,
			MaxAge:   maxAge,
			SameSite: http.SameSiteLaxMode,
		})
		if err = session.Save(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": 1, "message": "Save session error! "})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": 0})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Old Password Error!"})
	}
}

// UpdateProfile Update user profile
func UpdateProfile(c *gin.Context) {
	user := c.Value("user").(*models.User)
	email := c.PostForm("email")
	nickName := c.PostForm("nickname")
	g := c.PostForm("gender")
	bio := c.PostForm("bio")
	if g == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Parameters!"})
		return
	} else {
		gender, err := strconv.Atoi(g)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Parameters!"})
			return
		}
		err = service.UpdateProfile(user, email, nickName, gender, bio)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Parameters!"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"success": 0})
			return
		}
	}
}

func ChangeEncryptionAlgorithm(c *gin.Context) {
	user := c.Value("user").(*models.User)
	newAlgorithm := c.PostForm("algorithm")
	algo, err := strconv.Atoi(newAlgorithm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Parameters!"})
	}
	err = service.ChangeEncryptionAlgorithm(user, algo, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": 1, "message": "Server Error"})
	}
	session := sessions.Default(c)
	session.Delete("user")
	session.Delete("EncryptionKey")
	session.Delete("remember")
	session.Options(sessions.Options{
		Path:     "/api",
		HttpOnly: true,
		// MaxAge=0 will delete cookie after browser close
		MaxAge:   0,
		SameSite: http.SameSiteLaxMode,
	})
	if err = session.Save(); err != nil {
		c.JSON(http.StatusOK, gin.H{"success": 1, "message": "Save session error"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": 0})
	}
}
