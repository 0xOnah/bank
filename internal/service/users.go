package service

import (
	"context"

	"github.com/onahvictor/bank/internal/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg entity.User) (*entity.User, error)
	GetUser(ctx context.Context, username string) (*entity.User, error)
}

type userService struct {
	userRepo UserRepository
}

func (us *userService) CreateUser(ctx context.Context, username, email, password, fullname string) (entity.User, error) {
	user, err := entity.NewUser(username, password, fullname, email)
	if err != nil {
		return entity.User{}, err
	}
	createdUser, err := us.userRepo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, nil
	}
	return *createdUser, nil
}