package common_models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Post struct {
	ID        int32
	Content   string
	UserID    pgtype.Int4
	CreatedAt time.Time
	UpdatesAt time.Time
}
