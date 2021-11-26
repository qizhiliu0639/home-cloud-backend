package utils

import (
	"github.com/shiena/ansicolor"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var globalLogger *logrus.Logger
var loggerOnce sync.Once

func buildLogger() {
	globalLogger = logrus.New()
	globalLogger.SetLevel(logrus.InfoLevel)
	globalLogger.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	globalLogger.SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
}

// GetLogger return the logger instance
func GetLogger() *logrus.Logger {
	loggerOnce.Do(buildLogger)
	return globalLogger
}
