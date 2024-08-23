package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func getFormatter(isJSON bool) logrus.Formatter {
	if isJSON {
		return &logrus.JSONFormatter{
			PrettyPrint: true,
		}
	}
	return &logrus.TextFormatter{
		FullTimestamp:          false,
		DisableLevelTruncation: true,
	}
}

func Get() *logrus.Logger {
	if logger != nil {
		return logger
	}
	stdOutHandler := os.Stdout

	logger := &logrus.Logger{
		Out:       stdOutHandler,
		Formatter: getFormatter(true),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}

	return logger
}
