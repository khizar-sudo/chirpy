package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/khizar-sudo/chirpy/internal/auth"
	"github.com/khizar-sudo/chirpy/internal/utils"
)

type loginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
}

type loginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func login(cfg *apiConfig) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		body := loginRequest{}

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&body); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
			return
		}

		if body.Email == "" || body.Password == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Email and password are required", nil)
			return
		}

		user, err := cfg.db.GetUser(req.Context(), body.Email)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			} else {
				utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
			}
			return
		}

		ok, err := auth.CheckPasswordHash(body.Password, user.HashedPassword)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error verifying password", err)
			return
		}

		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password", nil)
			return
		}

		var expiresIn time.Duration
		if body.ExpiresInSeconds == nil || *body.ExpiresInSeconds > 3600 {
			expiresIn = time.Hour
		} else {
			expiresIn = time.Second * time.Duration(*body.ExpiresInSeconds)
		}

		token, err := auth.MakeJWT(user.ID, cfg.tokenSecret, expiresIn)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Could not make token", err)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, loginResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		})
	}
}
