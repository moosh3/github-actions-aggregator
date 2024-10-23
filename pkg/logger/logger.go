package logger

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)
	log.SetFormatter(&logrus.JSONFormatter{})
}

func GetLogger() *logrus.Logger {
	return log
}
