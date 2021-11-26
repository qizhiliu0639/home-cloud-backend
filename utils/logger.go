package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
	"sync"
)

var globalLogger *logrus.Logger
var loggerOnce sync.Once

func buildLogger() {
	globalLogger = logrus.New()
	globalLogger.SetLevel(logrus.InfoLevel)
	globalLogger.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\r\n",
	})
}

// GetLogger return the logger instance
func GetLogger() *logrus.Logger {
	loggerOnce.Do(buildLogger)
	return globalLogger
}
