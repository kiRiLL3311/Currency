package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/kiRiLL3311/Currency/backend/internal/middleware"
	"github.com/kiRiLL3311/Currency/backend/internal/services"
)

type RateHandler struct {
	Service *services.RateService
}

// GetRates godoc
// @Summary List exchange rates
// @Description Returns all stored USD-based exchange rates, including previous_rate and change_percent when available.
// @Tags Rates
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Rate
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /rates [get]
func (h *RateHandler) GetRates(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Request from user %d", userID)

	rates, err := h.Service.GetRates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}

// Convert godoc
// @Summary Convert currency
// @Description Converts an amount between two currencies using stored USD-based rates.
// @Tags Converter
// @Produce json
// @Security BearerAuth
// @Param from query string true "Source currency code" example(USD)
// @Param to query string true "Target currency code" example(EUR)
// @Param amount query number true "Amount to convert" example(1000)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {string} string "Invalid amount"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /convert [get]
func (h *RateHandler) Convert(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Conversion requested by user %d", userID)

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

// SyncRates godoc
// @Summary Sync rates from upstream
// @Description Public endpoint that refreshes rates from Open ER API and previous-day sources.
// @Tags Rates
// @Produce plain
// @Param base query string false "Base currency" default(USD) example(USD)
// @Success 200 {string} string "Rates synchronized successfully"
// @Failure 500 {string} string "Internal server error"
// @Router /rates/sync [get]
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
