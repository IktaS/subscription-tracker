package service

import (
	"context"

	"github.com/IktaS/subscription-tracker/entity"
)

type Notifier interface {
	NotifySubsription(ctx context.Context, subscription entity.Subscription) error
}
