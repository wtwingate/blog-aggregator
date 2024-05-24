package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(r *http.Request) (string, error) {
	apiKey := r.Header.Get("Authorization")
	apiKey = strings.TrimPrefix(apiKey, "ApiKey ")
	if len(apiKey) == 0 {
		return "", errors.New("missing API key")
	}
	return apiKey, nil
}
