package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func ValidateID(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.PostForm(key)
		if len(id) > 0 {
			vID, err := uuid.Parse(id)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Request"})
				c.Abort()
				return
			} else {
				c.Set(key, vID)
				c.Next()
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Request"})
			c.Abort()
			return
		}
	}
}
