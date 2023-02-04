package entity

import (
	"errors"
	"strconv"
	"time"
)

type User struct {
	ID     string
	Payday Payday
	Role   string
}

func (u User) IsValid() bool {
	return u.ID != ""
}

type Subscription struct {
	ID            string
	User          User
	Title         string
	PaymentMethod string
	Amount        Amount
	LastPaidDate  time.Time
	NextPaidDate  time.Time
	Duration      SubDuration
}

func (s Subscription) IsValid() bool {
	if !s.User.IsValid() {
		return false
	}
	if s.Title == "" {
		return false
	}
	if s.PaymentMethod == "" {
		return false
	}
	if !s.Amount.IsValid() {
		return false
	}
	if s.LastPaidDate.IsZero() {
		return false
	}
	if !s.Duration.IsValid() {
		return false
	}
	return true
}

func (s Subscription) GetNextPaymentDate() time.Time {
	switch s.Duration.Unit {
	case DurationUnitYear:
		return s.LastPaidDate.AddDate(s.Duration.Value, 0, 0)
	case DurationUnitMonth:
		return s.LastPaidDate.AddDate(0, s.Duration.Value, 0)
	case DurationUnitDay:
		return s.LastPaidDate.AddDate(0, 0, s.Duration.Value)
	default:
		return s.LastPaidDate.Add(s.Duration.ToTimeDuration())
	}
}

type Amount struct {
	Value    float64
	Currency string
}

func (a Amount) IsValid() bool {
	return len(a.Currency) == 3
}

type SubDurationUnit string

const (
	DurationUnitYear   SubDurationUnit = "year"
	DurationUnitMonth  SubDurationUnit = "month"
	DurationUnitDay    SubDurationUnit = "day"
	DurationUnitHour   SubDurationUnit = "hour"
	DurationUnitMinute SubDurationUnit = "minute"
)

func (s SubDuration) ToTimeDuration() time.Duration {
	switch s.Unit {
	case DurationUnitYear:
		return time.Hour * 24 * 365 * time.Duration(s.Value)
	case DurationUnitMonth:
		return time.Hour * 24 * 30 * time.Duration(s.Value)
	case DurationUnitDay:
		return time.Hour * 24 * time.Duration(s.Value)
	case DurationUnitHour:
		return time.Hour * time.Duration(s.Value)
	case DurationUnitMinute:
		return time.Minute * time.Duration(s.Value)
	default:
		return 0
	}
}

func StringToSubDurationUnit(s string) (SubDurationUnit, error) {
	var err error
	switch s {
	case string(DurationUnitYear):
		return DurationUnitYear, err
	case string(DurationUnitMonth):
		return DurationUnitMonth, err
	case string(DurationUnitDay):
		return DurationUnitDay, err
	case string(DurationUnitHour):
		return DurationUnitHour, err
	case string(DurationUnitMinute):
		return DurationUnitMinute, err
	default:
		return "", errors.New("unknown sub duration unit")
	}
}

type SubDuration struct {
	Value int
	Unit  SubDurationUnit
}

func (s SubDuration) IsValid() bool {
	_, err := StringToSubDurationUnit(string(s.Unit))
	return err == nil
}

type Payday string

const (
	End Payday = "end"
)

func StringToPayday(s string) (Payday, error) {
	if s == string(End) {
		return End, nil
	}
	_, err := strconv.Atoi(s)
	if err != nil {
		return "", err
	}
	return Payday(s), nil
}

func (p Payday) GetPaydayFromTo(t time.Time) (time.Time, time.Time, error) {
	if p == "" {
		return time.Time{}, time.Time{}, errors.New("invalid payday")
	}
	if p == End {
		firstOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		firstOfMonth = firstOfMonth.AddDate(0, 0, -1)
		return firstOfMonth, lastOfMonth, nil
	}
	val, err := strconv.Atoi(string(p))
	if err != nil {
		return time.Time{}, time.Time{}, nil
	}
	if t.Day() < val {
		firstOfPD := time.Date(t.Year(), t.Month()-1, val, 0, 0, 0, 0, t.Location())
		lastOfPD := firstOfPD.AddDate(0, 1, 0)
		return firstOfPD, lastOfPD, nil
	}
	firstOfPD := time.Date(t.Year(), t.Month(), val, 0, 0, 0, 0, t.Location())
	lastOfPD := firstOfPD.AddDate(0, 1, 0)
	return firstOfPD, lastOfPD, nil
}
