package sqlite

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/IktaS/subscription-tracker/entity"
)

const (
	getDefaultLogChannel = "select log_channel from configuration where user_id = $1;"
	setDefaultLogChannel = "insert or replace into configuration(user_id, log_channel) values($1, $2);"
)

func (s *SQLiteStore) GetDefaultLogChannel(ctx context.Context, user entity.User) (string, error) {
	var logChannel string
	err := s.db.QueryRowContext(ctx, getDefaultLogChannel, user.ID).Scan(&logChannel)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return logChannel, nil
}

func (s *SQLiteStore) SetDefaultLogChannel(ctx context.Context, user entity.User, logChannel string) error {
	stmt, err := s.db.PrepareContext(ctx, setDefaultLogChannel)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.ID, logChannel)
	if err != nil {
		return err
	}
	return nil
}
