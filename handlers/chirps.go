package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/khizar-sudo/chirpy/internal/database"
	"github.com/khizar-sudo/chirpy/internal/utils"
)

type chirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func createChirp(cfg *apiConfig) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		body := chirpRequest{}

		if err := decoder.Decode(&body); err != nil {
			utils.RespondWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s\n", err), err)
			return
		}

		if len(body.Body) > 140 {
			utils.RespondWithError(w, 400, "Chirp is too long", nil)
			return
		}

		words := strings.Split(body.Body, " ")
		for i, word := range words {
			w := strings.ToLower(word)
			if w == "kerfuffle" || w == "sharbert" || w == "fornax" {
				words[i] = "****"
			}
		}

		chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
			Body:   strings.Join(words, " "),
			UserID: body.UserID,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		}

		utils.RespondWithJSON(w, http.StatusCreated, chirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
}
