package common_models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Post struct {
	ID        int32
	Content   string
	UserID    pgtype.Int4
	CreatedAt time.Time
	UpdatesAt time.Time
}
