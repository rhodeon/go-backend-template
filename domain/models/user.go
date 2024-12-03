package models

import (
	"time"

	"github.com/rhodeon/go-backend-template/repositories/database/implementation/users"
)

type User struct {
	ID        int32
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) FromDbUser(dbUser users.User) User {
	return User(dbUser)
}
