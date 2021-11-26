package controllers

import (
	"github.com/gin-gonic/gin"
	"home-cloud/models"
	"net/http"
)

// GetUserStatus Check if logged in and get user status
func GetUserStatus(c *gin.Context) {
	user := c.Value("user").(*models.User)
	encryption := false
	if user.Encryption > 0 {
		encryption = true
	}
	c.JSON(http.StatusOK, gin.H{
		"success":         0,
		"username":        user.Username,
		"status":          user.Status,
		"email":           user.Email,
		"nickname":        user.Nickname,
		"gender":          user.Gender,
		"bio":             user.Bio,
		"account_salt":    user.AccountSalt,
		"encryption":      encryption,
		"encryption_algo": user.Encryption,
	})
}
