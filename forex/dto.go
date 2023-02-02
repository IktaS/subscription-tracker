package forex

import "errors"

type getCurrencyRateRequest struct {
	Base   string `json:"base"`
	Output string `json:"output"`
}

type getCurrencyRateResponse struct {
	Valid   bool               `json:"valid"`
	Updated int64              `json:"updated"`
	Base    string             `json:"base"`
	Rates   map[string]float64 `json:"rates"`
}

var ErrCurrencyNotFound = errors.New("currency not found")
