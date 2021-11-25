package middleware

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"strings"
)

func exist(fSys fs.FS, path string) bool {
	f, err := fs.Stat(fSys, path)
	if err != nil {
		return false
	}
	//Skip directory
	if f.IsDir() {
		return false
	}
	return true
}

// FrontendFileHandler return the static file in go embed if it is not under /api
func FrontendFileHandler(staticFS embed.FS, root string) gin.HandlerFunc {
	fSys, err := fs.Sub(staticFS, root)
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(fSys))
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/api") {
			c.Next()
			return
		}
		if !strings.HasPrefix(p, "/") {
			c.JSON(http.StatusBadRequest, gin.H{"success": 1})
			c.Abort()
			return
		}
		//Rewrite 404 page to index.html
		//If it is a directory, also return index.html
		if p == "/" || strings.HasSuffix(p, "/") || !exist(fSys, strings.TrimPrefix(p, "/")) {
			c.FileFromFS("/", http.FS(fSys))
		} else if strings.HasSuffix(p, "/index.html") {
			//Trim the index.html at the end to prevent redirection in ServeHTTP
			c.FileFromFS(strings.TrimSuffix(p, "index.html"), http.FS(fSys))
		} else {
			fileServer.ServeHTTP(c.Writer, c.Request)
		}
		c.Abort()
		return
	}
}
