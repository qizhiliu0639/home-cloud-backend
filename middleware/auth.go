package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"home-cloud/models"
	"home-cloud/utils"
	"net/http"
)

func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username, ok := session.Get("user").(string)
		if !ok || len(username) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": 1, "message": "You have not logged in!"})
			return
		}
		user, err := models.GetUserByUsername(username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": 1, "message": "You have not logged in!"})
			return
		}
		utils.GetLogger().Info("User " + user.Username + " request comes")
		c.Set("user", user)
		c.Next()
	}
}
