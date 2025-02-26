package app

import "github.com/saygik/go-url-shortener/internal/lib/logger"

func (a *App) initLogger(env string) {
	l := logger.InitLogger(env)
	l.Info("------------------Starting programm-------------")
	a.log = l
}
