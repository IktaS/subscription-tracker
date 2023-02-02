package service

type Forex interface {
	ToIDR(currency string, value float64) (float64, error)
}
