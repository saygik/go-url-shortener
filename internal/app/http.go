package app

import (
	"github.com/saygik/go-url-shortener/internal/handler/http"
	"github.com/saygik/go-url-shortener/internal/handler/http/api"
)

func (a *App) StartHTTPServer() error {

	s := http.NewServer(a.Config.Env, a.log)

	api.NewHandler(s.Rtr, a.C.GetUseCase(), a.log)

	err := s.Start(a.Config.Port)
	if err != nil {
		return err
	}

	return nil
}
