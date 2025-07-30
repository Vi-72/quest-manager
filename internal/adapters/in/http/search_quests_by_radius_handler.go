package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"quest-manager/internal/core/application/usecases/queries"
	"quest-manager/internal/core/domain/model/kernel"
)

// SearchQuestsByRadiusHTTPHandler handles GET /api/v1/quests/search-radius requests.
type SearchQuestsByRadiusHTTPHandler struct {
	queryHandler queries.SearchQuestsByRadiusQueryHandler
}

// NewSearchQuestsByRadiusHTTPHandler creates a new instance of SearchQuestsByRadiusHTTPHandler.
func NewSearchQuestsByRadiusHTTPHandler(handler queries.SearchQuestsByRadiusQueryHandler) *SearchQuestsByRadiusHTTPHandler {
	return &SearchQuestsByRadiusHTTPHandler{queryHandler: handler}
}

// ServeHTTP processes the HTTP request for searching quests within a radius.
func (h *SearchQuestsByRadiusHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract and validate query parameters
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	radiusStr := r.URL.Query().Get("radius_km")

	if latStr == "" || lonStr == "" || radiusStr == "" {
		http.Error(w, "lat, lon, and radius_km are required", http.StatusBadRequest)
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid lat value", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "invalid lon value", http.StatusBadRequest)
		return
	}

	radiusKm, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil || radiusKm <= 0 {
		http.Error(w, "invalid radius_km value", http.StatusBadRequest)
		return
	}

	center, err := kernel.NewGeoCoordinate(lat, lon)
	if err != nil {
		http.Error(w, "invalid coordinates: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Prepare query
	query := queries.SearchQuestsByRadiusQuery{
		Center:   center,
		RadiusKm: radiusKm,
	}

	result, err := h.queryHandler.Handle(r.Context(), query)
	if err != nil {
		http.Error(w, "failed to search quests: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return results
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result.Quests)
}
