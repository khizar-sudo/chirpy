package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/khizar-sudo/chirpy/internal/auth"
	"github.com/khizar-sudo/chirpy/internal/database"
	"github.com/khizar-sudo/chirpy/internal/utils"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&body); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		if body.Email == "" || body.Password == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Email and Password required", nil)
			return
		}

		hashPassword, err := auth.HashPassword(body.Password)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			return
		}

		user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
			Email:          body.Email,
			HashedPassword: hashPassword,
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Could not create user", err)
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
