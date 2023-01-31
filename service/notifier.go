package service

import "context"

type Notifier interface {
	PingUser(ctx context.Context, userID string) error
}
