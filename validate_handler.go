package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/khizar-sudo/chirpy/utils"
)

type validateBody struct {
	Body string `json:"body"`
}

type validateResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	body := validateBody{}
	w.Header().Set("Content-Type", "application/json")

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

	w.WriteHeader(http.StatusOK)
	data, err := json.Marshal(validateResponse{
		CleanedBody: strings.Join(words, " "),
	})
	if err != nil {
		utils.RespondWithError(w, 500, err.Error(), err)
	}
	w.Write(data)

}
