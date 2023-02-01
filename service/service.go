package service

import (
	"context"

	"github.com/IktaS/subscription-tracker/entity"
)

type Service interface {
	NotifyUserSubscription(ctx context.Context, user entity.User, sub entity.Subscription) error
}

type service struct {
	notifier Notifier
	logger   Logger
}

func NewService(notifier Notifier, logger Logger) *service {
	return &service{
		notifier: notifier,
		logger:   logger,
	}
}

func (s *service) NotifyUserSubscription(ctx context.Context, user entity.User, sub entity.Subscription) error {
	s.logger.Info(ctx, "info1")
	s.logger.Error(ctx, "info2")
	s.logger.Warning(ctx, "info3")
	sub.User = user
	return s.notifier.NotifySubsription(ctx, sub)
}
