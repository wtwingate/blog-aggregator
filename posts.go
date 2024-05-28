package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/wtwingate/blog-aggregator/internal/database"
)

func (cfg *apiConfig) handlerGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	params := r.URL.Query()
	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		limit = 20
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		errMsg := fmt.Sprintf("could not get posts: %v", err)
		respondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	respondWithJSON(w, http.StatusOK, toPostSlice(posts))
}

func (cfg *apiConfig) createNewPost(item Item, feedID uuid.UUID) error {
	ctx := context.Background()

	pubDate, err := parsePubDate(item.PubDate)
	if err != nil {
		return err
	}

	post, err := cfg.DB.CreatePost(ctx, database.CreatePostParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Title:     item.Title,
		Url:       item.Link,
		Description: sql.NullString{
			String: item.Description,
			Valid:  len(item.Description) > 0,
		},
		PublishedAt: pubDate,
		FeedID:      feedID,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "unique_violation" {
				return nil
			}
		}
		return err
	}

	log.Printf("created new post: %v\n", post.Title)
	return nil
}

func parsePubDate(pubDate string) (time.Time, error) {
	layouts := []string{
		time.RFC3339, time.RFC3339Nano, time.RFC1123Z, time.RFC1123, time.RFC822, time.RFC822Z,
	}
	for _, layout := range layouts {
		time, err := time.Parse(layout, pubDate)
		if err == nil {
			return time, nil
		}
	}
	return time.Now(), fmt.Errorf("could not parse time stamp %v", pubDate)
}
