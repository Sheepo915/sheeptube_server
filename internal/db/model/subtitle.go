package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type Subtitle struct {
	ID         int32        `db:"id"`
	SubtitleID uuid.UUID    `db:"subtitle_id"`
	Language   language.Tag `db:"language"`
	FilePath   string       `db:"file_path"`
	BelongTo   int32        `db:"belong_to"` // The video that own the subtitle
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at"`
}
