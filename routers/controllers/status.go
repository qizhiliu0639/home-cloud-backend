package controllers

import (
	"github.com/gin-gonic/gin"
	"home-cloud/models"
	"net/http"
)

// GetUserStatus Check if logged in and get user status
func GetUserStatus(c *gin.Context) {
	user := c.Value("user").(*models.User)
	c.JSON(http.StatusOK, gin.H{"username": user.Username, "status": user.Status})
}
