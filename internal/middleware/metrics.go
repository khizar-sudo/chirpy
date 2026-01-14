package middleware

import (
	"net/http"

	"github.com/khizar-sudo/chirpy/internal/config"
)

func MetricsInc(cfg *config.ApiConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg.FileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
	}
}
