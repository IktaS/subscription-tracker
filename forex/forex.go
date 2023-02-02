package forex

import (
	"log"
	"sync"
	"time"
)

type currencyRates struct {
	rates  map[string]float64
	expiry time.Time
}

type Forex struct {
	key        string
	logger     *log.Logger
	expiryTime int
	rates      map[string]currencyRates
	lock       *sync.RWMutex
}

// expiry in secods
func NewForexService(key string, expiry int, logger *log.Logger) *Forex {
	srv := &Forex{
		key:        key,
		expiryTime: expiry,
		logger:     logger,
		rates:      make(map[string]currencyRates),
		lock:       &sync.RWMutex{},
	}
	return srv
}

func (f *Forex) ToIDR(base string, value float64) (float64, error) {
	// var err error
	// var rates currencyRates
	// f.lock.RLock()
	// rates, ok := f.rates["IDR"]
	// f.lock.RUnlock()
	// if !ok || time.Now().After(rates.expiry) {
	// 	rates, err = f.callGetRatesEndpoint("IDR")
	// 	if err != nil {
	// 		return -1, err
	// 	}
	// 	f.lock.Lock()
	// 	f.rates["IDR"] = rates
	// 	f.lock.Unlock()
	// }
	// rate, ok := rates.rates[base]
	// if !ok {
	// 	return -1, ErrCurrencyNotFound
	// }
	// return value * rate, nil
	var rate float64
	switch base {
	case "USD":
		rate = 17000
	case "JPY":
		rate = 120
	default:
		return -1, ErrCurrencyNotFound
	}
	return rate * value, nil
}
