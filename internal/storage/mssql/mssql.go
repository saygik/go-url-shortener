package mssql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/saygik/go-url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

type ConnectionParameters struct {
	Server   string
	Database string
	User     string
	Password string
}

func New(cp ConnectionParameters) (*Storage, error) {
	const op = "storage.mssql.New"

	dbinfo := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", cp.Server, cp.Database, cp.User, cp.Password)
	var err error
	db, err := sql.Open("mssql", dbinfo)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// db, err := sql.Open("sqlite3", storagePath)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }
	//defer db.Close()

	// stmt, err := db.Prepare(`
	// CREATE TABLE IF NOT EXISTS url(
	// 	id INTEGER PRIMARY KEY,
	// 	alias TEXT NOT NULL UNIQUE,
	// 	url TEXT NOT NULL);
	// CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	// `)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	// _, err = stmt.Exec()
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
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

func (s *Storage) GetURL(alias string) (string, error) {
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
