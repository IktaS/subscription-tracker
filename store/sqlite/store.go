package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	filename string
	db       *sql.DB
}

func NewSQLiteStore(filename string) (*SQLiteStore, error) {
	var err error
	s := &SQLiteStore{
		filename: filename,
	}
	s.db, err = sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return s, err
}

func (s *SQLiteStore) Shutdown() {
	s.db.Close()
}
