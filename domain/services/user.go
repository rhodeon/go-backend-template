package services

import (
	"context"
	"strings"

	domainerrors "github.com/rhodeon/go-backend-template/domain/errors"
	"github.com/rhodeon/go-backend-template/domain/models"
	"github.com/rhodeon/go-backend-template/repositories"
	"github.com/rhodeon/go-backend-template/repositories/database/postgres"
	pgusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type User struct {
	*service
}

var userService *User

func newUser(repos *repositories.Repositories) *User {
	userService = &User{newService(repos)}
	return userService
}

func (u User) Create(ctx context.Context, dbTx postgres.Transaction, user models.User) (models.User, error) {
	_, err := u.repos.Database.Users.Create(ctx, dbTx, pgusers.CreateParams{
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), postgres.UniqueUsersEmail):
			return models.User{}, domainerrors.NewDuplicateDataError("user", "email", user.Email)

		case strings.Contains(err.Error(), postgres.UniqueUsersUsername):
			return models.User{}, domainerrors.NewDuplicateDataError("user", "username", user.Username)

		default:
			return models.User{}, errors.Wrap(err, "unable to create user")
		}
	}

	return models.User{}, nil
}

func (u User) GetById(ctx context.Context, dbTx postgres.Transaction, userId int32) (models.User, error) {
	dbUser, err := u.repos.Database.Users.GetById(ctx, dbTx, int64(userId))
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
