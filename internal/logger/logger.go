package logger

import "github.com/sirupsen/logrus"

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	Logger.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	Logger.SetFormatter(formatter)
}
