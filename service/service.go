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
	forex    Forex
}

func NewService(notifier Notifier, logger Logger, forex Forex) *service {
	return &service{
		notifier: notifier,
		logger:   logger,
		forex:    forex,
	}
}

func (s *service) NotifyUserSubscription(ctx context.Context, user entity.User, sub entity.Subscription) error {
	var err error
	sub.User = user
	sub, err = s.convertCurrencyIfNotIDR(ctx, sub)
	if err != nil {
		s.logger.Error(ctx, "failed to convert to IDR", "subscription", sub, "error", err)
		return err
	}
	s.logger.Info(ctx, "hey hey hey")
	return s.notifier.NotifySubsription(ctx, sub)
}

func (s *service) convertCurrencyIfNotIDR(ctx context.Context, sub entity.Subscription) (entity.Subscription, error) {
	var err error
	if sub.Amount.Currency != "IDR" {
		newSub := sub
		newSub.Amount.Currency = "IDR"
		newSub.Amount.Value, err = s.forex.ToIDR(sub.Amount.Currency, sub.Amount.Value)
		return newSub, err
	}
	return sub, err
}
