package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Video struct {
	ID           int32          `db:"id"`
	VideoID      uuid.UUID      `db:"video_id"`
	Name         string         `db:"name"`
	Description  sql.NullString `db:"description"`
	Source       string         `db:"source"`
	PostedBy     int32          `db:"posted_by"` // Posted by channel
	Likes        int32          `db:"likes"`
	Views        int32          `db:"views"`
	Saved        int32          `db:"saved"`
	CommentCount int64          `db:"comment_count"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}
