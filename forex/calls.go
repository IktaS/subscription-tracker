package forex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	baseURL       = "https://currencyapi.net/api/v1/"
	defaultOutput = "JSON"
)

func (f *Forex) callGetRatesEndpoint(base string) (currencyRates, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(baseURL+"/rates?key=%s&base=%s&output=%s", f.key, base, defaultOutput),
		nil,
	)
	if err != nil {
		f.logger.Printf("error creating request for GetRatesEndpoint: %v", err)
		return currencyRates{}, err
	}
	r, err := client.Do(req)
	if err != nil {
		f.logger.Printf("error doing request for GetRatesEndpoint: %v", err)
		return currencyRates{}, err
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return currencyRates{}, err
	}
	str := string(body)
	fmt.Println(str)
	dec := json.NewDecoder(r.Body)
	var resp getCurrencyRateResponse
	err = dec.Decode(&resp)
	if err != nil {
		f.logger.Printf("error decoding response for GetRatesEndpoint: %v", err)
		return currencyRates{}, err
	}
	ret := currencyRates{
		rates:  resp.Rates,
		expiry: time.Now().Add(time.Second * time.Duration(f.expiryTime)),
	}
	return ret, nil
}
