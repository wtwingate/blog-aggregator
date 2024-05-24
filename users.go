package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wtwingate/blog-aggregator/internal/database"
)

type usersPost struct {
	Name string `json:"name"`
}

func (cfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	request := usersPost{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ctx := context.Background()

	response, err := cfg.DB.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      request.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusOK, response)
}
