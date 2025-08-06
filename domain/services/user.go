package services

import (
	"context"
	"strings"

	domainerrors "github.com/rhodeon/go-backend-template/domain/errors"
	"github.com/rhodeon/go-backend-template/domain/models"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/repositories"
	"github.com/rhodeon/go-backend-template/repositories/database/postgres"
	dbusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"
	"github.com/rhodeon/go-backend-template/utils/typeutils"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	*service
}

var userService *User

func newUser(repos *repositories.Repositories) *User {
	userService = &User{newService(repos)}
	return userService
}

func (u *User) Create(ctx context.Context, dbTx *database.Tx, user models.User) (models.User, error) {
	hashedPassword, err := u.hashPassword(user.Password)
	if err != nil {
		return models.User{}, errors.Wrap(err, "hashing password")
	}

	createdUser, err := u.repos.Database.Users.Create(ctx, dbTx, dbusers.CreateParams{
		Username:       user.Username,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		PhoneNumber:    postgres.NewPgxText(typeutils.Ptr(user.PhoneNumber)),
		HashedPassword: hashedPassword,
	})
	if err != nil {
		switch {
		case strings.Contains(err.Error(), postgres.UniqueUsersEmail):
			return models.User{}, domainerrors.NewDuplicateDataError("user", "email", user.Email)

		case strings.Contains(err.Error(), postgres.UniqueUsersUsername):
			return models.User{}, domainerrors.NewDuplicateDataError("user", "username", user.Username)

		default:
			return models.User{}, errors.Wrap(err, "creating user in database")
		}
	}

	return models.NewUser.FromDbUser(createdUser), nil
}

func (u *User) GetById(ctx context.Context, dbTx *database.Tx, userId int64) (models.User, error) {
	dbUser, err := u.repos.Database.Users.GetById(ctx, dbTx, userId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return models.User{}, domainerrors.NewRecordNotFoundErr("user")

		default:
			return models.User{}, errors.Wrap(err, "getting user by id from database")
		}
	}

	return models.User{}.FromDbUser(dbUser), nil
}

func (u *User) hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.Wrap(err, "generating password with bcrypt")
	}

	return string(passwordHash), nil
}
