package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
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

// GetPreviousDayRates returns base→target rates for a recent prior day.
// Tries up to 7 calendar days back so weekends/holidays still resolve.
func GetPreviousDayRates(base string) (map[string]float64, error) {
	baseLower := strings.ToLower(base)

	var lastErr error
	day := time.Now().UTC().AddDate(0, 0, -1)

	for i := 0; i < 7; i++ {
		date := day.AddDate(0, 0, -i).Format("2006-01-02")
		url := fmt.Sprintf(
			"https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@%s/v1/currencies/%s.min.json",
			date,
			baseLower,
		)

		resp, err := http.Get(url)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("historical rates %s: status %d", date, resp.StatusCode)
			continue
		}

		var raw map[string]json.RawMessage
		err = json.NewDecoder(resp.Body).Decode(&raw)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		payload, ok := raw[baseLower]
		if !ok {
			lastErr = fmt.Errorf("historical rates %s: missing %s key", date, baseLower)
			continue
		}

		var src map[string]float64
		if err := json.Unmarshal(payload, &src); err != nil {
			lastErr = err
			continue
		}
		if len(src) == 0 {
			lastErr = fmt.Errorf("historical rates %s: empty payload", date)
			continue
		}

		out := make(map[string]float64, len(src))
		for code, rate := range src {
			out[strings.ToUpper(code)] = rate
		}
		return out, nil
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("historical rates unavailable")
	}
	return nil, lastErr
}
