package services

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rhodeon/go-backend-template/models"
	"github.com/rhodeon/go-backend-template/repositories"
	"github.com/rhodeon/go-backend-template/repositories/database"
	"github.com/rhodeon/go-backend-template/repositories/database/implementation/users"
	"log/slog"
	"strings"
)

type User struct {
	*service
}

var userService *User

func newUser(repos *repositories.Repositories, logger *slog.Logger) *User {
	userService = &User{newService(repos, logger)}
	return userService
}

func (u User) Create(ctx context.Context, dbTx database.Transaction, user models.User) (models.User, error) {
	dbCreatedUser, err := u.repos.Database.Users.Create(ctx, dbTx, users.CreateParams{
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), database.UniqueUsersEmail):
		case strings.Contains(err.Error(), database.UniqueUsersUsername):
		default:
			return models.User{}, errors.Wrap(err, "unable to create user")
		}
	}

	return models.User{}.FromDbUser(dbCreatedUser), nil
}
