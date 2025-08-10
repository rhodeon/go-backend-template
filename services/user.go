package services

import (
	"context"
	"strings"

	"github.com/rhodeon/go-backend-template/domain"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/repositories"
	"github.com/rhodeon/go-backend-template/repositories/database/postgres"
	dbusers "github.com/rhodeon/go-backend-template/repositories/database/postgres/sqlcgen/users"
	"github.com/rhodeon/go-backend-template/utils/typeutils"

	"github.com/go-errors/errors"
	"github.com/jackc/pgx/v5"
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

func (u *User) Create(ctx context.Context, dbTx *database.Tx, user domain.User) (domain.User, error) {
	hashedPassword, err := u.hashPassword(user.Password)
	if err != nil {
		return domain.User{}, errors.Errorf("hashing password: %w", err)
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
			return domain.User{}, domain.UserErrDuplicateEmail(user.Email)

		case strings.Contains(err.Error(), postgres.UniqueUsersUsername):
			return domain.User{}, domain.UserErrDuplicateUsername(user.Username)

		default:
			return domain.User{}, errors.Errorf("creating user in database: %w", err)
		}
	}

	return domain.NewUser.FromDbUser(createdUser), nil
}

func (u *User) GetById(ctx context.Context, dbTx *database.Tx, userId int64) (domain.User, error) {
	dbUser, err := u.repos.Database.Users.GetById(ctx, dbTx, userId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.User{}, domain.UserErrNotFound("id", userId)

		default:
			return domain.User{}, errors.Errorf("getting user with id %q from database: %w", userId, err)
		}
	}

	return domain.NewUser.FromDbUser(dbUser), nil
}

func (u *User) hashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.Errorf("generating password with bcrypt: %w", err)
	}

	return string(passwordHash), nil
}
