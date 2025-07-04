package service

import (
	"context"
	"log"

	"github.com/onahvictor/bank/internal/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg entity.User) (*entity.User, error)
	GetUser(ctx context.Context, username string) (*entity.User, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(ur UserRepository) *userService {
	return &userService{
		userRepo: ur,
	}
}

func (us *userService) CreateUser(ctx context.Context, username, email, password, fullname string) (entity.User, error) {
	user, err := entity.NewUser(username, password, fullname, email)
	if err != nil {
		log.Print(err)
		return entity.User{}, NewAppError(ErrBadRequest, "failed to create user", err)
	}
	createdUser, err := us.userRepo.CreateUser(ctx, user)
	if err != nil {
		return entity.User{}, NewAppError(ErrInternal, "internal server error", err)
	}
	return *createdUser, nil
}
