package middleware

import (
	"github.com/gin-gonic/gin"
	"home-cloud/utils"
	"time"
)

// LoggerToFile Log to file
func LoggerToFile() gin.HandlerFunc {

	logger:=utils.GetLogger()

	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logger.Infof("| %3d | %13v |%13v | %15s | %s | %s |",
			statusCode,
			startTime,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}