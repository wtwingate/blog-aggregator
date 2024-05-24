package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/wtwingate/blog-aggregator/internal/database"
)

type authHandlerFunc func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		if len(apiKey) == 0 {
			respondWithError(w, http.StatusBadRequest, "missing API key")
			return
		}

		split := strings.Fields(apiKey)
		if len(split) < 2 || split[0] != "ApiKey" {
			respondWithError(w, http.StatusBadRequest, "malformed API key")
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), split[1])
		if err != nil {
			errMsg := fmt.Sprintf("invalid API key: %v", err)
			respondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}

		handler(w, r, user)
	}
}
