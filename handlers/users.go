package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/khizar-sudo/chirpy/utils"
)

type userRequest struct {
	Email string `json:"email"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func createUser(cfg *apiConfig) func(w http.ResponseWriter, req *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body := userRequest{}
		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&body); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
			return
		}

		user, err := cfg.db.CreateUser(req.Context(), body.Email)
		if err != nil {
			utils.RespondWithError(w, 500, "Could not create user", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, userResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		})
	})
}
