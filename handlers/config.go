package handlers

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/khizar-sudo/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (cfg *apiConfig) getMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
}
