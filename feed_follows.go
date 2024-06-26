package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wtwingate/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := fmt.Sprintf("could not decode parameters: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	feedFollows, err := cfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		errMsg := fmt.Sprintf("could not follow feed: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	respondWithJSON(w, http.StatusCreated, toFeedFollow(feedFollows))
}

func (cfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		errMsg := fmt.Sprintf("could not get feed follows: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
	respondWithJSON(w, http.StatusOK, toFeedFollowSlice(feedFollows))
}

func (cfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	ffidString := r.PathValue("feedFollowID")
	if len(ffidString) == 0 {
		respondWithError(w, http.StatusBadRequest, "missing feed follow ID")
		return
	}

	feedFollowID, err := uuid.Parse(ffidString)
	if err != nil {
		errMsg := fmt.Sprintf("could not parse feed follow ID: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	err = cfg.DB.DeleteFeedFollows(r.Context(), feedFollowID)
	if err != nil {
		errMsg := fmt.Sprintf("could not delete feed follow: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
}
