package sqlite

import (
	"context"
	"database/sql"
)

const (
	getDefaultLogChannel = "select log_channel from configuration;"
	setDefaultLogChannel = "insert or replace into configuration(log_channel) values(?);"
)

func (s *SQLiteStore) GetDefaultLogChannel(ctx context.Context) (string, error) {
	var logChannel string
	err := s.db.QueryRow(getDefaultLogChannel).Scan(&logChannel)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return logChannel, nil
}

func (s *SQLiteStore) SetDefaultLogChannel(ctx context.Context, logChannel string) error {
	stmt, err := s.db.PrepareContext(ctx, setDefaultLogChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, logChannel)
	if err != nil {
		return err
	}
	return nil
}
