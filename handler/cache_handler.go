package handler

import (
	"encoding/json"
	"net/http"
	"task-manager-api/dto"
	"task-manager-api/usecase"
)

func RegisterCacheRoutes(mux *http.ServeMux, cache *usecase.CacheService) {
	mux.HandleFunc("GET /cache/stats", func(w http.ResponseWriter, r *http.Request) {
		getCacheStats(w, r, cache)
	})
	mux.HandleFunc("DELETE /cache", func(w http.ResponseWriter, r *http.Request) {
		clearCache(w, r, cache)
	})
}

func getCacheStats(w http.ResponseWriter, r *http.Request, cache *usecase.CacheService) {
	stats := cache.GetStats()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func clearCache(w http.ResponseWriter, r *http.Request, cache *usecase.CacheService) {
	cache.Clear()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.SuccessResponse{
		Message: "Cache cleared successfully",
	})
}
