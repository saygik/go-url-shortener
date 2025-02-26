package logger

import (
	"fmt"
	"os"

	"github.com/saygik/go-url-shortener/internal/lib/constants"
	"github.com/sirupsen/logrus"
)

type Log struct {
}

func InitLogger(env string) *logrus.Logger {

	var log = logrus.New()

	src, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("err", err)
	}
	log.Out = src

	switch env {
	case constants.EnvLocal:
		customFormatter := new(logrus.TextFormatter)
		customFormatter.TimestampFormat = "2006-01-02 15:04:05"
		customFormatter.FullTimestamp = true
		log.Formatter = customFormatter
	case constants.EnvProd:
		customFormatter2 := new(logrus.JSONFormatter)
		customFormatter2.TimestampFormat = "2006-01-02 15:04:05"
		log.Formatter = customFormatter2
		log.Debug("logger debug mode enabled")
	}

	return log
}

// func Err22(err error) slog.Attr {
// 	return slog.Attr{
// 		Key:   "error",
// 		Value: slog.StringValue(err.Error()),
// 	}
// }
