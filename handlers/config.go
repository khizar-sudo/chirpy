package handlers

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/khizar-sudo/chirpy/internal/database"
	"github.com/khizar-sudo/chirpy/internal/utils"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func (cfg *apiConfig) getMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, req *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		utils.RespondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	cfg.fileserverHits.Store(0)

	err := cfg.db.DeleteAllUsers(req.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Could not delete users", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
