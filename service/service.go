package service

import (
	"context"
	"fmt"
	"time"

	"github.com/IktaS/subscription-tracker/entity"
	"github.com/google/uuid"
	"github.com/procyon-projects/chrono"
)

type Service interface {
	NotifyUserSubscription(ctx context.Context, user entity.User, sub entity.Subscription) error
	LoadSubscriptions(ctx context.Context) error
	GetAllSubscriptionForUser(ctx context.Context, user entity.User) ([]entity.Subscription, error)
	GetAllSubscriptionForUserUntilPayday(ctx context.Context, user entity.User) ([]entity.Subscription, error)
	GetAllSubscriptionInPaydayCycle(ctx context.Context, user entity.User, t time.Time) ([]entity.Subscription, error)
	SetPaydayTime(ctx context.Context, user entity.User) error
	SetSubscription(ctx context.Context, sub entity.Subscription) error
}

type service struct {
	store     Store
	notifier  Notifier
	logger    Logger
	forex     Forex
	scheduler chrono.TaskScheduler
}

func NewService(store Store, notifier Notifier, logger Logger, forex Forex) *service {
	return &service{
		store:     store,
		notifier:  notifier,
		logger:    logger,
		forex:     forex,
		scheduler: chrono.NewDefaultTaskScheduler(),
	}
}

func (s *service) Shutdown() bool {
	return <-s.scheduler.Shutdown()
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
	subs, err := s.store.LoadSubscriptions(ctx)
	if err != nil {
		return err
	}
	for _, sub := range subs {
		now := time.Now().UTC()
		for sub.NextPaidDate.Before(now) {
			sub.LastPaidDate = sub.NextPaidDate
			sub.NextPaidDate = sub.GetNextPaymentDate()
		}
		err := s.store.SetSubscription(ctx, sub)
		if err != nil {
			s.logger.Error(ctx, "failed to set subscription when loading", "error", err)
		}
		err = s.scheduleSubscriptionNotification(ctx, sub)
		if err != nil {
			s.logger.Error(ctx, "failed to set scheduler when loading", "error", err)
		}
	}
	return nil
}

func (s *service) GetAllSubscriptionForUser(ctx context.Context, user entity.User) ([]entity.Subscription, error) {
	var err error
	subs, err := s.store.GetAllSubscriptionsForUser(ctx, user)
	if err != nil {
		return nil, err
	}
	sbs := s.calculateSubTotals(ctx, subs)
	return sbs, nil
}

func (s *service) GetAllSubscriptionInPaydayCycle(ctx context.Context, user entity.User, cycle time.Time) ([]entity.Subscription, error) {
	var err error
	user.Payday, err = s.store.GetPaydayTime(ctx, user)
	if err != nil {
		return nil, err
	}
	subs, err := s.store.GetAllSubscriptionsForUserInPaydayCycle(ctx, user, cycle)
	if err != nil {
		return nil, err
	}
	sbs := s.calculateSubTotals(ctx, subs)
	return sbs, nil
}

func (s *service) calculateSubTotals(ctx context.Context, subs []entity.Subscription) []entity.Subscription {
	var sbs []entity.Subscription
	totals := make(map[string]float64)
	for _, v := range subs {
		idr, err := s.forex.ToIDR(v.Amount.Currency, v.Amount.Value)
		if err != nil {
			s.logger.Error(ctx, "failed to convert to IDR", "error", err)
			continue
		}
		totals[v.PaymentMethod] += idr
		sbs = append(sbs, entity.Subscription{
			ID:            v.ID,
			Title:         v.Title,
			User:          v.User,
			PaymentMethod: v.PaymentMethod,
			Amount: entity.Amount{
				Currency: "IDR",
				Value:    idr,
			},
			LastPaidDate: v.LastPaidDate,
			NextPaidDate: v.NextPaidDate,
			Duration:     v.Duration,
		})
	}
	for k, v := range totals {
		sbs = append(sbs, entity.Subscription{
			Title:         fmt.Sprintf("Total %s", k),
			PaymentMethod: k,
			Amount: entity.Amount{
				Currency: "IDR",
				Value:    v,
			},
		})
	}
	return sbs
}

func (s *service) GetAllSubscriptionForUserUntilPayday(ctx context.Context, user entity.User) ([]entity.Subscription, error) {
	var err error
	user.Payday, err = s.store.GetPaydayTime(ctx, user)
	if err != nil {
		return nil, err
	}
	subs, err := s.store.GetAllSubscriptionsForUserUntilPayday(ctx, user)
	if err != nil {
		return nil, err
	}
	sbs := s.calculateSubTotals(ctx, subs)
	return sbs, nil
}

func (s *service) SetPaydayTime(ctx context.Context, user entity.User) error {
	return s.store.SetPaydayTime(ctx, user)
}

func (s *service) SetSubscription(ctx context.Context, sub entity.Subscription) error {
	if sub.ID == "" {
		sub.ID = uuid.NewString()
	}
	sub.LastPaidDate = sub.LastPaidDate.UTC()
	sub.NextPaidDate = sub.GetNextPaymentDate()
	err := s.store.SetSubscription(ctx, sub)
	if err != nil {
		return err
	}
	return s.scheduleSubscriptionNotification(ctx, sub)
}

func (s *service) scheduleSubscriptionNotification(ctx context.Context, sub entity.Subscription) error {
	_, err := s.scheduler.Schedule(func(ctx context.Context) {
		err := s.NotifyUserSubscription(ctx, sub.User, sub)
		if err != nil {
			s.logger.Error(ctx, "failed to notify user subscription", "error", err)
		}
		sub.LastPaidDate = time.Now().UTC()
		sub.NextPaidDate = sub.GetNextPaymentDate()
		err = s.SetSubscription(ctx, sub)
		if err != nil {
			s.logger.Error(ctx, "failed to set user subscription", "error", err)
		}
	}, chrono.WithTime(sub.NextPaidDate.Local()))
	return err
}
