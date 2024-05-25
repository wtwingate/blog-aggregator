package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wtwingate/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerCreateFeeds(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := fmt.Sprintf("could not decode parameters: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		errMsg := fmt.Sprintf("could not create feed: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	respondWithJSON(w, http.StatusOK, dbFeedToFeed(feed))
}

func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		errMsg := fmt.Sprintf("could not get feeds: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	respondWithJSON(w, http.StatusOK, dbFeedsToFeeds(feeds))
}
