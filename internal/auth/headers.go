package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	token := strings.Split(headers.Get("Authorization"), " ")[1]

	if token == "" {
		return token, errors.New("No token found")
	}

	return token, nil
}