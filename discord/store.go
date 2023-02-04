package discord

import (
	"context"

	"github.com/IktaS/subscription-tracker/entity"
)

type Store interface {
	GetDefaultLogChannel(ctx context.Context, user entity.User) (string, error)
	SetDefaultLogChannel(ctx context.Context, user entity.User, logChannel string) error
}
