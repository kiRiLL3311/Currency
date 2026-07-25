package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kiRiLL3311/Currency/news/internal/middleware"
	"github.com/kiRiLL3311/Currency/news/internal/services"
)

type NewsHandler struct {
	Service *services.NewsService
}

// List godoc
// @Summary List news articles
// @Description Returns stored headlines for a region, optionally filtered by keyword.
// @Tags News
// @Produce json
// @Security BearerAuth
// @Param region query string false "Region code" Enums(US,EU,GB,JP,AU) example(US)
// @Param keyword query string false "Keyword used when pulling / filtering" example(euro)
// @Success 200 {array} models.Article
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /news [get]
func (h *NewsHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	region := r.URL.Query().Get("region")
	keyword := r.URL.Query().Get("keyword")
	articles, err := h.Service.List(userID, region, keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

// Sync godoc
// @Summary Sync news from upstream
// @Description Pulls articles from NewsAPI and/or RSS for a region (and keyword). Public endpoint.
// @Tags News
// @Produce plain
// @Param region query string false "Region code (omit with no keyword to sync all)" Enums(US,EU,GB,JP,AU) example(US)
// @Param keyword query string false "Keyword included in upstream pull" example(yen)
// @Success 200 {string} string "News synchronized successfully"
// @Failure 500 {string} string "Internal server error"
// @Router /news/sync [get]
func (h *NewsHandler) Sync(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	keyword := r.URL.Query().Get("keyword")

	if region == "" && keyword == "" {
		h.Service.SyncAll()
		w.Write([]byte("News synchronized for all regions"))
		return
	}

	if region == "" {
		region = "US"
	}

	if err := h.Service.SyncRegion(region, keyword); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("News synchronized successfully"))
}
