package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/wtwingate/blog-aggregator/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func toUser(user database.User) User {
	return User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		ApiKey:    user.ApiKey,
	}
}

type Feed struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Url           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
}

func toFeed(feed database.Feed) Feed {
	return Feed{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		Url:           feed.Url,
		UserID:        feed.UserID,
		LastFetchedAt: &feed.LastFetchedAt.Time,
	}
}

func toFeedSlice(feedSlice []database.Feed) []Feed {
	feeds := []Feed{}
	for _, f := range feedSlice {
		feeds = append(feeds, toFeed(f))
	}
	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func toFeedFollow(feedFollows database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        feedFollows.ID,
		CreatedAt: feedFollows.CreatedAt,
		UpdatedAt: feedFollows.CreatedAt,
		UserID:    feedFollows.UserID,
		FeedID:    feedFollows.FeedID,
	}
}

func toFeedFollowSlice(feedFollowSlice []database.FeedFollow) []FeedFollow {
	feedFollows := []FeedFollow{}
	for _, ff := range feedFollowSlice {
		feedFollows = append(feedFollows, toFeedFollow(ff))
	}
	return feedFollows
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func toPost(post database.Post) Post {
	return Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: &post.Description.String,
		PublishedAt: post.PublishedAt,
		FeedID:      post.FeedID,
	}
}

func toPostSlice(postSlice []database.Post) []Post {
	posts := []Post{}
	for _, post := range postSlice {
		posts = append(posts, toPost(post))
	}
	return posts
}
