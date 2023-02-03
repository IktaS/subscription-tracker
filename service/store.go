package service

import (
	"context"

	"github.com/IktaS/subscription-tracker/entity"
)

type Store interface {
	LoadSubscriptions(ctx context.Context) ([]entity.Subscription, error)
	GetAllSubscriptionsForUser(ctx context.Context, user entity.User) ([]entity.Subscription, error)
	SetSubscription(ctx context.Context, sub entity.Subscription) error
	SetPaydayTime(ctx context.Context, user entity.User) error
	GetPaydayTime(ctx context.Context, user entity.User) (entity.Payday, error)
}
