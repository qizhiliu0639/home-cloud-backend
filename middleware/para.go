package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ValidateDir() gin.HandlerFunc {
	return func(c *gin.Context) {
		dir := c.PostForm("dir")
		if len(dir) > 0 && strings.HasPrefix(dir, "/") && (!strings.HasSuffix(dir[1:], "/")) {
			//filter root slash
			path := dir[1:]
			var paths []string
			//not a root path
			if len(path) > 0 {
				pathsTmp := strings.Split(path, "/")
				for _, p := range pathsTmp {
					//filter invalid path
					if len(p) < 1 || p == "." || p == ".." {
						c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Path"})
						c.Abort()
						return
					} else {
						paths = append(paths, p)
					}
				}
			}
			c.Set("vDir", paths)
			c.Next()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"success": 1, "message": "Invalid Path"})
			c.Abort()
		}
	}
}
