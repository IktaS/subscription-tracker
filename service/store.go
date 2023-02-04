package service

import (
	"context"
	"time"

	"github.com/IktaS/subscription-tracker/entity"
)

type Store interface {
	LoadSubscriptions(ctx context.Context) ([]entity.Subscription, error)
	GetAllSubscriptionsForUser(ctx context.Context, user entity.User) ([]entity.Subscription, error)
	GetAllSubscriptionsForUserUntilPayday(ctx context.Context, user entity.User) ([]entity.Subscription, error)
	GetAllSubscriptionsForUserInPaydayCycle(ctx context.Context, user entity.User, cycle time.Time) ([]entity.Subscription, error)
	SetSubscription(ctx context.Context, sub entity.Subscription) error
	SetPaydayTime(ctx context.Context, user entity.User) error
	GetPaydayTime(ctx context.Context, user entity.User) (entity.Payday, error)
}
