package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/wtwingate/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := fmt.Sprintf("could not decode parameters: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		errMsg := fmt.Sprintf("could not create user: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	respondWithJSON(w, http.StatusOK, dbUserToUser(user))
}

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authorization")
	apiKey = strings.TrimPrefix(apiKey, "ApiKey ")
	if len(apiKey) == 0 {
		respondWithError(w, http.StatusUnauthorized, "missing API key")
		return
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		errMsg := fmt.Sprintf("could not get user: %v\n", err)
		respondWithError(w, http.StatusUnauthorized, errMsg)
		return
	}

	respondWithJSON(w, http.StatusOK, dbUserToUser(user))
}
