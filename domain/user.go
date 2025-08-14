package domain

import (
	"time"

	dbusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"
)

var NewUser = new(User)

type User struct {
	Id          int64
	Email       string
	Username    string
	FirstName   string
	LastName    string
	PhoneNumber string
	IsVerified  bool
	// Password is the hashed password for an already existing user.
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) FromDbUser(dbUser dbusers.User) User {
	return User{
		Id:          dbUser.Id,
		Email:       dbUser.Email,
		Username:    dbUser.Username,
		FirstName:   dbUser.FirstName,
		LastName:    dbUser.LastName,
		PhoneNumber: dbUser.PhoneNumber.String,
		Password:    dbUser.HashedPassword,
		IsVerified:  dbUser.IsVerified,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
	}
}
