package app

import (
	"database/sql"
	"fmt"

	"github.com/saygik/go-url-shortener/internal/config"
	"github.com/saygik/go-url-shortener/internal/storage/mssql"
)

func (a *App) newMsSQLConnect(cfg config.DBConfig) (*sql.DB, error) {
	db, err := mssql.NewConnection(fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", cfg.DBServer, cfg.DBName, cfg.DBUser, cfg.DBPassword))
	if err != nil {
		return nil, err
	}
	return db, nil
}
