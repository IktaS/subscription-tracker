package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/IktaS/subscription-tracker/entity"
	"github.com/google/uuid"
)

const (
	loadAllSubscription = `SELECT 
			id,
			user_id,
			title,
			payment_method,
			amount_currency,
			amount_value,
			last_paid,
			next_paid,
			duration_value,
			duration_unit 
		FROM subscription;`
	getAllSubscriptionForUser = `SELECT
			id,
			user_id,
			title,
			payment_method,
			amount_currency,
			amount_value,
			last_paid,
			next_paid,
			duration_value,
			duration_unit 
		FROM subscription WHERE user_id=$1;`
	getAllSubscriptionForUserInPaydayCycle = `SELECT
			id,
			user_id,
			title,
			payment_method,
			amount_currency,
			amount_value,
			last_paid,
			next_paid,
			duration_value,
			duration_unit 
		FROM subscription WHERE user_id=$1 and next_paid BETWEEN $2 and $3;`
	getAllSubscriptionForUserUntilPayday = `SELECT
			id,
			user_id,
			title,
			payment_method,
			amount_currency,
			amount_value,
			last_paid,
			next_paid,
			duration_value,
			duration_unit 
		FROM subscription WHERE user_id=$1 and next_paid <= $2;`
	setSubscription = `INSERT OR REPLACE INTO subscription(
		id, 
		user_id, 
		title, 
		payment_method, 
		amount_currency,
		amount_value,
		last_paid,
		next_paid,
		duration_value,
		duration_unit
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	setPaydayTime = `INSERT OR REPLACE INTO user(id, payday_time) values($1, $2);`
	getPaydayTime = `SELECT payday_time FROM user WHERE id = $1`
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
			&sub.ID,
			&sub.User.ID,
			&sub.Title,
			&sub.PaymentMethod,
			&sub.Amount.Currency,
			&sub.Amount.Value,
			&sub.LastPaidDate,
			&sub.NextPaidDate,
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
			&sub.ID,
			&sub.User.ID,
			&sub.Title,
			&sub.PaymentMethod,
			&sub.Amount.Currency,
			&sub.Amount.Value,
			&sub.LastPaidDate,
			&sub.NextPaidDate,
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

func (s *SQLiteStore) GetAllSubscriptionsForUserInPaydayCycle(ctx context.Context, user entity.User, cycle time.Time) ([]entity.Subscription, error) {
	prevPD, nextPD, err := user.Payday.GetPaydayFromTo(cycle)
	if err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, getAllSubscriptionForUserInPaydayCycle, user.ID, prevPD.Format(time.RFC3339), nextPD.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []entity.Subscription
	for rows.Next() {
		var sub entity.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.User.ID,
			&sub.Title,
			&sub.PaymentMethod,
			&sub.Amount.Currency,
			&sub.Amount.Value,
			&sub.LastPaidDate,
			&sub.NextPaidDate,
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

func (s *SQLiteStore) GetAllSubscriptionsForUserUntilPayday(ctx context.Context, user entity.User) ([]entity.Subscription, error) {
	_, nextPD, err := user.Payday.GetPaydayFromTo(time.Now())
	if err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(ctx, getAllSubscriptionForUserInPaydayCycle, user.ID, nextPD.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subs []entity.Subscription
	for rows.Next() {
		var sub entity.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.User.ID,
			&sub.Title,
			&sub.PaymentMethod,
			&sub.Amount.Currency,
			&sub.Amount.Value,
			&sub.LastPaidDate,
			&sub.NextPaidDate,
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
	if sub.ID == "" {
		sub.ID = uuid.NewString()
	}
	stmt, err := s.db.PrepareContext(
		ctx,
		setSubscription,
	)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(
		ctx,
		sub.ID,
		sub.User.ID,
		sub.Title,
		sub.PaymentMethod,
		sub.Amount.Currency,
		sub.Amount.Value,
		sub.LastPaidDate,
		sub.NextPaidDate,
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
		user.ID,
		user.Payday,
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
