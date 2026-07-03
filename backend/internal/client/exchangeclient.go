package client

import (
	"encoding/json"
	"net/http"
)

type ExchangeResponse struct {
	BaseCode string             `json:"base_code"`
	Rates    map[string]float64 `json:"rates"`
}

func GetRates(base string) (*ExchangeResponse, error) {
	resp, err := http.Get("https://open.er-api.com/v6/latest/" + base)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ExchangeResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
