package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	InitLogger()
}

func InitLogger() {
	log = logrus.New()

	// Set log format
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set output to stdout or a file
	log.SetOutput(os.Stdout)

	// Optional: Set log level
	log.SetLevel(logrus.InfoLevel)
}

func GetLogger() *logrus.Logger {
	return log
}
