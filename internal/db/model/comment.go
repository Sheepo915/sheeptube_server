package model

import (
	"time"

	"github.com/google/uuid"
)

type Comments struct {
	ID             int32     `db:"id"`
	CommentID      uuid.UUID `db:"comment_id"`
	Content        string    `db:"content"`
	CommentedBy    int32     `db:"commented_by"`
	Likes          int32     `db:"likes"`
	ParentComment  int32     `db:"parent_comment"`
	VideoCommented int32     `db:"video_commented"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
