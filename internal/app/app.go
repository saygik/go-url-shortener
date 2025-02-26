package app

import (
	"github.com/saygik/go-url-shortener/internal/config"
	"github.com/sirupsen/logrus"
)

type App struct {
	Config config.Config
	C      *Container
	log    *logrus.Logger
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	a := &App{
		Config: cfg,
	}
	a.initLogger(cfg.Env)

	cnn, err := a.newMsSQLConnect(cfg.DB)
	if err != nil {
		return nil, err
	}
	//	storage, _ :=
	a.C = NewContainer(cnn)
	return a, nil
}
