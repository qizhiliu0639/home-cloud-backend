package utils

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var globalLogger *logrus.Logger
var loggerOnce sync.Once
var logPath = GetConfig().LogFilePath

func buildLogger() {
	globalLogger = logrus.New()
	//logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err != nil {
	//	globalLogger.Errorf("[Error] Open log file error: %s. Use default stderr", err)
	//}
	//globalLogger.Out = logFile
	globalLogger.SetLevel(logrus.InfoLevel)
	globalLogger.SetFormatter(&logrus.TextFormatter{})
}

func GetLogger() *logrus.Logger {
	loggerOnce.Do(buildLogger)
	return globalLogger
}
