package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/onahvictor/bank/internal/config"
	"github.com/onahvictor/bank/internal/db/repo"
	"github.com/onahvictor/bank/internal/entity"
	"github.com/onahvictor/bank/internal/sdk/auth"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg entity.User) (*entity.User, error)
	GetUser(ctx context.Context, username string) (*entity.User, error)
}

type userService struct {
	userRepo UserRepository
	token    auth.Auntenticator
	config   *config.Config
}

func NewUserService(ur UserRepository, token auth.Auntenticator, config config.Config) *userService {
	return &userService{
		userRepo: ur,
		token:    token,
		config:   &config,
	}
}

func (us *userService) CreateUser(ctx context.Context, username, email, password, fullname string) (entity.User, error) {
	user, err := entity.NewUser(username, password, fullname, email)
	if err != nil {
		return entity.User{}, NewAppError(ErrBadRequest, "failed to create user", err)
	}
	createdUser, err := us.userRepo.CreateUser(ctx, user)
	if err != nil {
		errvalue, ok := err.(*pq.Error)
		if ok {
			switch {
			case strings.Contains(errvalue.Error(), "duplicate"):
				return entity.User{}, NewAppError(ErrBadRequest, "this user already exist", err)
			}
		}

		return entity.User{}, NewAppError(ErrInternal, "internal server error", err)
	}
	return *createdUser, nil
}

type AuthResult struct {
	Token string
	User  *entity.User
}

func (us *userService) Login(ctx context.Context, username, password string) (*AuthResult, error) {
	user, err := us.userRepo.GetUser(ctx, username)
	if err != nil {
		if err == repo.ErrUserNotFound {
			return nil, NewAppError(ErrBadRequest, "user not found", err)
		}
		return nil, NewAppError(ErrInternal, "internal server error:", err)
	}

	if !auth.ComparePassword([]byte(user.HashedPassword), password) {
		return nil, repo.ErrInvalidCredentials
	}

	token, err := us.token.GenerateToken(user.Username, us.config.ACCESS_TOKEN_DURATATION)
	if err != nil {
		return nil, fmt.Errorf("token generation failed: %w", err)
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil

}
