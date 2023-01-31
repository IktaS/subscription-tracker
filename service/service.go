package service

import "context"

type Service interface {
	PingUser(ctx context.Context, userID string) error
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

func (s *service) PingUser(ctx context.Context, userID string) error {
	return s.notifier.PingUser(ctx, userID)
}
