package middleware

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

func exist(fSys fs.FS, path string) bool {
	f, err := fs.Stat(fSys, path)
	if os.IsNotExist(err) {
		return false
	}
	if f.IsDir() {
		return false
	}
	return true
}

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
		//Prevent directory listing
		if strings.HasSuffix(p, "/") {
			p = path.Join(p, "index.html")
		}
		//Filter the first slash
		p = p[1:]
		//Rewrite 404 page to index.html
		if !exist(fSys, p) {
			//Do not use c.FileFromFS here, or will cause infinite redirection in index.html<=>./
			index, err := staticFS.ReadFile(path.Join(root, "index.html"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": 1})
				c.Abort()
				return
			} else {
				c.Header("Content-Type", "text/html")
				c.String(http.StatusOK, string(index))
				c.Abort()
				return
			}
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return

	}

}
