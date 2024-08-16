// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package posts

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Post struct {
	ID        int32       `db:"id"`
	Content   string      `db:"content"`
	UserID    pgtype.Int4 `db:"user_id"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatesAt time.Time   `db:"updates_at"`
}