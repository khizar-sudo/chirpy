package handlers

import (
	"fmt"
	"net/http"

	"github.com/khizar-sudo/chirpy/internal/config"
)

func GetMetrics(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.FileserverHits.Load())
	}
}

func ResetMetrics(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		cfg.FileserverHits.Store(0)
	}
}
