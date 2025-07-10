package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID          int32          `db:"id"`
	ChannelID   uuid.UUID      `db:"channel_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	OwnedBy     string         `db:"owned_by"`
	CreatedAt   time.Time      `db:"created_at"`
}
