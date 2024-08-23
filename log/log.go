package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func getFormatter(isJSON bool) logrus.Formatter {
	if isJSON {
		return &logrus.JSONFormatter{
			PrettyPrint:     false,
			TimestampFormat: "2006-01-02 15:04:05",
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				s := strings.Split(f.Function, ".")
				funcVal := s[len(s)-1]
				file := strings.Split(f.File, "/")
				folder := strings.Split(file[len(file)-1], "\\")
				fileVal := fmt.Sprintf("%s:%d", folder[len(folder)-1], f.Line)
				return funcVal, fileVal
			},
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
		Out:          stdOutHandler,
		Formatter:    getFormatter(true),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ReportCaller: true,
	}

	return logger
}
