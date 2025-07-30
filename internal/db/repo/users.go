package repo

import (
	"context"

	"github.com/0xOnah/bank/internal/db/sqlc"
	"github.com/0xOnah/bank/internal/entity"
)

type userRepo struct {
	db *sqlc.SQLStore
}

func NewUserRepo(db *sqlc.SQLStore) *userRepo {
	return &userRepo{
		db: db,
	}
}
func ToUser(u *sqlc.User) *entity.User {

	return &entity.User{
		Username:          u.Username,
		HashedPassword:    u.HashedPassword,
		Email:             entity.Email{Value: u.Email},
		FullName:          u.FullName,
		CreatedAt:         u.CreatedAt,
		PasswordChangedAt: u.PasswordChangedAt,
	}

}
func (ur *userRepo) CreateUser(ctx context.Context, arg entity.User) (*entity.User, error) {
	user, err := ur.db.CreateUser(ctx, sqlc.CreateUserParams{
		Username:       arg.Username,
		HashedPassword: arg.HashedPassword,
		FullName:       arg.FullName,
		Email:          arg.Email.String(),
	})
	if err != nil {
		return nil, err
	}

	return ToUser(user), nil

}

func (ur *userRepo) GetUser(ctx context.Context, username string) (*entity.User, error) {
	user, err := ur.db.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return ToUser(user), err
}
