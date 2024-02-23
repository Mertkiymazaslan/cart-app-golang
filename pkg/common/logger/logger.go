package logger

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func Initialize() (*logrus.Logger, error) {
	logger = logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	return logger, nil
}

func GetInstance() *logrus.Logger {
	return logger
}
