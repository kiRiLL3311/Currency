package models

type Rate struct {
	ID             int     `json:"id"`
	BaseCurrency   string  `json:"base_currency"`
	TargetCurrency string  `json:"target_currency"`
	Rate           float64 `json:"rate"`
}
