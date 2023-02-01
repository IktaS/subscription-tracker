package entity

import "time"

type User struct {
	ID string
}

type Subscription struct {
	User          User
	Title         string
	Description   string
	PaymentMethod string
	Amount        Amount
	LastPaidDate  time.Time
	Duration      SubDuration
}

type Amount struct {
	Value    float64
	Currency string
}

type SubDuration struct {
	Value    int
	Duration time.Duration
}
