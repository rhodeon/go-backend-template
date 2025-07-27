package commonmodels

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Post struct {
	Id        int32
	Content   string
	UserId    pgtype.Int4
	CreatedAt time.Time
	UpdatesAt time.Time
}
