package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xOnah/bank/internal/config"
	"github.com/0xOnah/bank/internal/db/repo"
	"github.com/0xOnah/bank/internal/entity"
	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/transport/sdk/errorutil"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg entity.User) (*entity.User, error)
	GetUser(ctx context.Context, username string) (*entity.User, error)
}
type SessionRepository interface {
	CreateSession(ctx context.Context, arg entity.Session) (*entity.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (*entity.Session, error)
}

type userService struct {
	userRepo    UserRepository
	token       auth.Authenticator
	config      *config.Config
	SessionRepo SessionRepository
}

func NewUserService(ur UserRepository, token auth.Authenticator, config config.Config, sr SessionRepository) *userService {
	return &userService{
		userRepo:    ur,
		token:       token,
		config:      &config,
		SessionRepo: sr,
	}
}

type CreateUserInput struct {
	Username string
	Password string
	Fullname string
	Email    string
}

func (us *userService) CreateUser(ctx context.Context, cr CreateUserInput) (entity.User, error) {
	user, err := entity.NewUser(cr.Username, cr.Password, cr.Fullname, cr.Email)
	fmt.Println(err)
	if err != nil {
		return entity.User{}, errorutil.NewAppError(errorutil.ErrBadRequest, "failed validation", err)
	}

	createdUser, err := us.userRepo.CreateUser(ctx, user)
	if err != nil {
		errvalue, ok := err.(*pq.Error)
		if ok {
			switch {
			case strings.Contains(errvalue.Error(), "duplicate"):
				return entity.User{}, errorutil.NewAppError(errorutil.ErrBadRequest, "this user already exist", err)
			}
		}
		return entity.User{}, errorutil.NewAppError(errorutil.ErrInternal, "internal server error", err)
	}
	return *createdUser, nil
}

type AuthResult struct {
	SessionID             uuid.UUID
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	User                  *entity.User
}

type Logininput struct {
	Username  string
	Password  string
	ClientIP  string
	UserAgent string
}

func (us *userService) Login(ctx context.Context, arg Logininput) (*AuthResult, error) {
	user, err := us.userRepo.GetUser(ctx, arg.Username)
	if err != nil {
		if err == repo.ErrUserNotFound {
			return nil, errorutil.NewAppError(errorutil.ErrBadRequest, "user not found", err)
		}
		return nil, errorutil.NewAppError(errorutil.ErrInternal, "internal server error:", err)
	}

	if !auth.ComparePassword([]byte(user.HashedPassword), arg.Password) {
		return nil, errorutil.NewAppError(errorutil.ErrUnauthorized, "wrong password", nil)
	}

	accessToken, accessPayload, err := us.token.GenerateToken(user.Username, us.config.ACCESS_TOKEN_DURATATION)
	if err != nil {
		return nil, errorutil.NewAppError(errorutil.ErrInternal, "token generation failed: %w", err)
	}

	refreshToken, refreshpayload, err := us.token.GenerateToken(user.Username, us.config.REFRESH_TOKEN_DURATION)
	if err != nil {
		return nil, errorutil.NewAppError(errorutil.ErrInternal, "token generation failed: %w", err)
	}

	session, err := us.SessionRepo.CreateSession(ctx, entity.Session{
		ID:           uuid.MustParse(refreshpayload.ID),
		Username:     user.Username,
		RefreshToken: refreshToken,
		ClientIp:     arg.ClientIP,
		UserAgent:    arg.UserAgent,
		IsBlocked:    false,
		ExpiresAt:    refreshpayload.ExpiresAt.Time,
	})
	if err != nil {
		return nil, errorutil.NewAppError(errorutil.ErrInternal, "failed to create session", err)
	}

	return &AuthResult{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt.Time,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshpayload.ExpiresAt.Time,
		User:                  user,
	}, nil

}

type RenewAccessToken struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
}

func (us *userService) RenewAccessToken(ctx context.Context, refreshToken string) (RenewAccessToken, error) {
	refreshPayload, err := us.token.VerifyToken(refreshToken)
	if err != nil {
		return RenewAccessToken{}, errorutil.NewAppError(errorutil.ErrUnauthorized, "session has expired", err)
	}

	session, err := us.SessionRepo.GetSession(ctx, uuid.MustParse(refreshPayload.ID))
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrSessionNotFound):
			return RenewAccessToken{}, errorutil.NewAppError(errorutil.ErrNotFound, "session not found", err)
		}
		return RenewAccessToken{}, errorutil.NewAppError(errorutil.ErrInternal, "internal server error", err)
	}

	if session.IsBlocked {
		return RenewAccessToken{}, errorutil.NewAppError(errorutil.ErrUnauthorized, "session is blocked", err)
	}
	if session.Username != refreshPayload.Username {
		return RenewAccessToken{}, errorutil.NewAppError(errorutil.ErrUnauthorized, "incorrect session user", err)
	}

	if session.RefreshToken != refreshToken {
		return RenewAccessToken{}, errorutil.NewAppError(errorutil.ErrUnauthorized, "mismatched session token", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return RenewAccessToken{}, errorutil.NewAppError(errorutil.ErrUnauthorized, "expired session", err)
	}

	accessToken, accessPayload, err := us.token.GenerateToken(refreshPayload.Username, us.config.ACCESS_TOKEN_DURATATION)
	if err != nil {
		return RenewAccessToken{}, errorutil.NewAppError(errorutil.ErrInternal, "internal server error", err)
	}

	return RenewAccessToken{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}, nil
}
