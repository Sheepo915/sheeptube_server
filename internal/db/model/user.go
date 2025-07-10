package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int32     `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
}
