package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/IktaS/subscription-tracker/entity"
	"github.com/google/uuid"
)

const (
	loadAllSubscription = `SELECT 
			user_id,
			title,
			payment_method,
			amount_currency,
			amount_value,
			last_paid,
			duration_value,
			duration_unit 
		FROM subscriptions;`
	getAllSubscriptionForUser = `SELECT 
			user_id,
			title,
			payment_method,
			amount_currency,
			amount_value,
			last_paid,
			duration_value,
			duration_unit 
		FROM subscriptions WHERE user_id=$1;`
	setSubscription = `INSERT INTO subscriptions(
		id, 
		user_id, 
		title, 
		payment_method, 
		amount_currency,
		amount_value,
		last_paid,
		duration_value,
		duration_unit
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	setPaydayTime = `INSERT OR REPLACE INTO user(id, payday_time) values($1, $2);`
	getPaydayTime = `SELECT payday_time WHERE id = $1`
)

func (s *SQLiteStore) LoadSubscriptions(ctx context.Context) ([]entity.Subscription, error) {
	rows, err := s.db.QueryContext(ctx, loadAllSubscription)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []entity.Subscription
	for rows.Next() {
		var sub entity.Subscription
		err := rows.Scan(
			&sub.User.ID,
			&sub.Title,
			&sub.PaymentMethod,
			&sub.Amount.Currency,
			&sub.Amount.Value,
			&sub.LastPaidDate,
			&sub.Duration.Value,
			&sub.Duration.Unit,
		)
		if err != nil {
			return nil, err
		}
		if !sub.IsValid() {
			return nil, errors.New("invalid subscription data")
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (s *SQLiteStore) GetAllSubscriptionsForUser(ctx context.Context, user entity.User) ([]entity.Subscription, error) {
	rows, err := s.db.QueryContext(ctx, getAllSubscriptionForUser, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []entity.Subscription
	for rows.Next() {
		var sub entity.Subscription
		err := rows.Scan(
			&sub.User.ID,
			&sub.Title,
			&sub.PaymentMethod,
			&sub.Amount.Currency,
			&sub.Amount.Value,
			&sub.LastPaidDate,
			&sub.Duration.Value,
			&sub.Duration.Unit,
		)
		if err != nil {
			return nil, err
		}
		if !sub.IsValid() {
			return nil, errors.New("invalid subscription data")
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (s *SQLiteStore) SetSubscription(ctx context.Context, sub entity.Subscription) error {
	stmt, err := s.db.PrepareContext(
		ctx,
		getAllSubscriptionForUser,
	)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		uuid.New(),
		sub.User.ID,
		sub.Title,
		sub.PaymentMethod,
		sub.Amount.Currency,
		sub.Amount.Value,
		sub.LastPaidDate,
		sub.Duration.Value,
		sub.Duration.Unit,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) SetPaydayTime(ctx context.Context, user entity.User) error {
	stmt, err := s.db.PrepareContext(
		ctx,
		setPaydayTime,
	)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		user.Payday,
		user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) GetPaydayTime(ctx context.Context, user entity.User) (entity.Payday, error) {
	var paydayTime string
	err := s.db.QueryRowContext(ctx, getPaydayTime, user.ID).Scan(&paydayTime)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return entity.StringToPayday(paydayTime)
}
