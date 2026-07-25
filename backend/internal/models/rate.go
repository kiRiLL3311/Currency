package models

type Rate struct {
	ID             int      `json:"id"`
	BaseCurrency   string   `json:"base_currency"`
	TargetCurrency string   `json:"target_currency"`
	Rate           float64  `json:"rate"`
	PreviousRate   *float64 `json:"previous_rate,omitempty"`
	ChangePercent  *float64 `json:"change_percent,omitempty"`
}
