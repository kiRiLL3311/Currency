package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kiRiLL3311/Currency/backend/internal/services"
)

type RateHandler struct {
	Service *services.RateService
}

// GetRates godoc
// @Summary Get all exchange rates
// @Description Returns all stored exchange rates
// @Tags Rates
// @Produce json
// @Success 200 {array} models.Rate
// @Router /rates [get]
func (h *RateHandler) GetRates(w http.ResponseWriter, r *http.Request) {
	rates, err := h.Service.GetRates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}

func (h *RateHandler) Convert(w http.ResponseWriter, r *http.Request) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	amountStr := r.URL.Query().Get("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	rate, converted, err := h.Service.Convert(from, to, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"from":      from,
		"to":        to,
		"amount":    amount,
		"rate":      rate,
		"converted": converted,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RateHandler) SyncRates(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("base")
	if base == "" {
		base = "USD"
	}

	err := h.Service.SyncRates(base)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Rates synchronized successfully"))
}
