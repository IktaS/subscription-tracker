package service

import (
	"context"
	"errors"

	"github.com/IktaS/subscription-tracker/entity"
)

type Service interface {
	NotifyUserSubscription(ctx context.Context, user entity.User, sub entity.Subscription) error
	LoadSubscriptions(ctx context.Context) error
	GetAllSubscriptionForUser(ctx context.Context, user entity.User) ([]entity.Subscription, error)
	GetAllSubscriptionUntilPayday(ctx context.Context, user entity.User) ([]entity.Subscription, error)
	SetPaydayTime(ctx context.Context, user entity.User) error
	SetSubscription(ctx context.Context, sub entity.Subscription) error
}

type service struct {
	store    Store
	notifier Notifier
	logger   Logger
	forex    Forex
}

func NewService(store Store, notifier Notifier, logger Logger, forex Forex) *service {
	return &service{
		store:    store,
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

func (s *service) LoadSubscriptions(ctx context.Context) error {
	return errors.New("unimplemented")
}

func (s *service) GetAllSubscriptionForUser(ctx context.Context, user entity.User) ([]entity.Subscription, error) {
	return nil, errors.New("unimplemented")
}

func (s *service) GetAllSubscriptionUntilPayday(ctx context.Context, user entity.User) ([]entity.Subscription, error) {
	return nil, errors.New("unimplemented")
}

func (s *service) SetPaydayTime(ctx context.Context, user entity.User) error {
	return s.store.SetPaydayTime(ctx, user)
}

func (s *service) SetSubscription(ctx context.Context, sub entity.Subscription) error {
	return s.store.SetSubscription(ctx, sub)
}
