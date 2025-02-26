package mssql

import (
	"database/sql"
	"errors"
	"fmt"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/saygik/go-url-shortener/internal/storage"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(conn *sql.DB) *Repository {
	return &Repository{
		db: conn,
	}
}

func (s *Repository) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO urlAliases (url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if mssqlErr, ok := err.(mssql.Error); ok && mssqlErr.Number == 2601 {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get count inserted rows: %w", op, err)
	}

	return id, nil
}

func (s *Repository) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT [url] FROM urlAliases WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string
	defer stmt.Close()
	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}
