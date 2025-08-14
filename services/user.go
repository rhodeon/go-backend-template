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
)

type User struct {
	*service
}

var userService *User

func newUser(repos *repositories.Repositories, cfg *Config) *User {
	userService = &User{newService(repos, cfg)}
	return userService
}

func (u *User) Create(ctx context.Context, dbTx *database.Tx, user domain.User) (domain.User, error) {
	hashedPassword, err := authService.hashPassword(user.Password)
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

func (u *User) SendVerificationEmail(ctx context.Context, user domain.User, otp string) error {
	if err := u.repos.Email.SendVerificationEmail(ctx, user.Email, otp); err != nil {
		return errors.Errorf("sending verification email: %w", err)
	}

	return nil
}

// Verify verifies a user's email using the given one-time password (OTP), and either updates their verification status or returns an error.
// Errors:
// - domain.ErrUserNotFound
// - domain.ErrUserAlreadyVerified
// - domain.ErrAuthInvalidOtp
func (u *User) Verify(ctx context.Context, dbTx *database.Tx, email string, otp string) error {
	dbUser, err := u.repos.Database.Users.GetByEmail(ctx, dbTx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return domain.UserErrNotFound("email", email)

		default:
			return errors.Errorf("getting user with email %q from database: %w", email, err)
		}
	}

	if dbUser.IsVerified {
		return domain.UserErrAlreadyVerified(dbUser.Id)
	}

	isValidToken, err := authService.ValidateOtp(ctx, otp, dbUser.Id)
	if err != nil {
		return errors.Errorf("validating otp: %w", err)
	}
	if !isValidToken {
		return domain.ErrAuthInvalidOtp
	}

	if err := u.repos.Database.Users.Verify(ctx, dbTx, dbUser.Id); err != nil {
		return errors.Errorf("verifying user with id %q in database: %w", dbUser.Id, err)
	}

	return nil
}
