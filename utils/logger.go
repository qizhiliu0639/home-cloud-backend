package utils

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var globalLogger *logrus.Logger
var loggerOnce sync.Once

func buildLogger() {
	globalLogger = logrus.New()
	globalLogger.SetLevel(logrus.InfoLevel)
	globalLogger.SetFormatter(&logrus.TextFormatter{})
}

func GetLogger() *logrus.Logger {
	loggerOnce.Do(buildLogger)
	return globalLogger
}
