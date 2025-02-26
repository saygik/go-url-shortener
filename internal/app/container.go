package app

import (
	"database/sql"

	"github.com/saygik/go-url-shortener/internal/storage/mssql"
	"github.com/saygik/go-url-shortener/internal/usecase"
)

type Container struct {
	mssql *sql.DB
}

func NewContainer(mssqlConnect *sql.DB) *Container {

	return &Container{
		mssql: mssqlConnect,
	}
}

func (c *Container) GetUseCase() *usecase.UseCase {

	return usecase.NewUseCase(c.getMssqlRepository())
}

func (c *Container) getMssql() *sql.DB {
	return c.mssql
}

func (c *Container) getMssqlRepository() usecase.Repository {
	return mssql.NewRepository(c.getMssql())
}
