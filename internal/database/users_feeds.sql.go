// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users_feeds.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUsersFeeds = `-- name: CreateUsersFeeds :one
INSERT INTO users_feeds (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, created_at, updated_at, user_id, feed_id
`

type CreateUsersFeedsParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

func (q *Queries) CreateUsersFeeds(ctx context.Context, arg CreateUsersFeedsParams) (UsersFeed, error) {
	row := q.db.QueryRowContext(ctx, createUsersFeeds,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	var i UsersFeed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.FeedID,
	)
	return i, err
}

const deleteUsersFeeds = `-- name: DeleteUsersFeeds :exec
DELETE FROM users_feeds WHERE id = $1
`

func (q *Queries) DeleteUsersFeeds(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUsersFeeds, id)
	return err
}