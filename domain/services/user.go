package services

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	domainerrors "github.com/rhodeon/go-backend-template/domain/errors"
	"github.com/rhodeon/go-backend-template/domain/models"
	"github.com/rhodeon/go-backend-template/repositories"
	"github.com/rhodeon/go-backend-template/repositories/database"
	"github.com/rhodeon/go-backend-template/repositories/database/implementation/users"
)

type User struct {
	*service
}

var userService *User

func newUser(repos *repositories.Repositories) *User {
	userService = &User{newService(repos)}
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
			return models.User{}, domainerrors.NewDuplicateDataError("user", "email", user.Email)

		case strings.Contains(err.Error(), database.UniqueUsersUsername):
			return models.User{}, domainerrors.NewDuplicateDataError("user", "username", user.Username)

		default:
			return models.User{}, errors.Wrap(err, "unable to create user")
		}
	}

	return models.User{}.FromDbUser(dbCreatedUser), nil
}

func (u User) GetById(ctx context.Context, dbTx database.Transaction, userId int32) (models.User, error) {
	dbUser, err := u.repos.Database.Users.GetById(ctx, dbTx, userId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return models.User{}, domainerrors.NewRecordNotFoundErr("user")

		default:
			return models.User{}, errors.Wrap(err, "unable to find user")
		}
	}

	return models.User{}.FromDbUser(dbUser), nil
}
