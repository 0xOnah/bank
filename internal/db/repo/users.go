package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xOnah/bank/internal/db/sqlc"
	"github.com/0xOnah/bank/internal/entity"
)

type UserRepo struct {
	db *sqlc.SQLStore
}

func NewUserRepo(db *sqlc.SQLStore) *UserRepo {
	return &UserRepo{
		db: db,
	}
}
func ToUser(u *sqlc.User) (*entity.User, error) {
	email, err := entity.NewEmail(u.Email)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		Username:          u.Username,
		HashedPassword:    u.HashedPassword,
		Email:             email,
		FullName:          u.FullName,
		CreatedAt:         u.CreatedAt,
		PasswordChangedAt: u.PasswordChangedAt,
	}, nil
}
func (ur *UserRepo) CreateUser(ctx context.Context, arg entity.User) (*entity.User, error) {
	user, err := ur.db.CreateUser(ctx, sqlc.CreateUserParams{
		Username:       arg.Username,
		HashedPassword: arg.HashedPassword,
		FullName:       arg.FullName,
		Email:          arg.Email.String(),
	})
	if err != nil {
		return nil, err
	}

	return ToUser(user)

}

func (ur *UserRepo) GetUser(ctx context.Context, username string) (*entity.User, error) {
	user, err := ur.db.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return ToUser(user)
}

type UpdateUserParams struct {
	FullName          *string
	HashedPassword    *string
	Email             *string
	Username          string
	PasswordChangedAt *time.Time
}

func deref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func (ur *UserRepo) UpdateUser(ctx context.Context, arg UpdateUserParams) (*entity.User, error) {
	user, err := ur.db.UpdateUser(ctx, sqlc.UpdateUserParams{
		FullName: sql.NullString{
			String: deref(arg.FullName),
			Valid:  arg.FullName != nil,
		},
		HashedPassword: sql.NullString{
			String: deref(arg.HashedPassword),
			Valid:  arg.HashedPassword != nil,
		},
		Email: sql.NullString{
			String: deref(arg.Email),
			Valid:  arg.Email != nil,
		},
		PasswordChangedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: arg.HashedPassword != nil,
		},
		Username: arg.Username,
	})
	if err != nil {
		return nil, err
	}
	return ToUser(user)
}
