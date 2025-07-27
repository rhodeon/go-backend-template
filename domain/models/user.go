package models

import (
	"time"

	pgusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"
)

type User struct {
	Id        int32
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) FromDbUser(dbUser pgusers.User) User {
	return User(dbUser)
}
