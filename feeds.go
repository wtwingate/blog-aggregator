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

	feedFollows, err := cfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    feed.UserID,
		FeedID:    feed.ID,
	})
	if err != nil {
		errMsg := fmt.Sprintf("could not follow feed: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	type FeedFeedFollow struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}

	FFF := FeedFeedFollow{
		Feed:       toFeed(feed),
		FeedFollow: toFeedFollow(feedFollows),
	}

	respondWithJSON(w, http.StatusCreated, FFF)
}

func (cfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		errMsg := fmt.Sprintf("could not get feeds: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	respondWithJSON(w, http.StatusOK, toFeedSlice(feeds))
}
