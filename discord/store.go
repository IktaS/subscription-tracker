package discord

import "context"

type Store interface {
	GetDefaultLogChannel(ctx context.Context) (string, error)
	SetDefaultLogChannel(ctx context.Context, logChannel string) error
}
