package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/onahvictor/bank/internal/config"
	"github.com/onahvictor/bank/internal/db/repo"
	"github.com/onahvictor/bank/internal/entity"
	"github.com/onahvictor/bank/internal/sdk/auth"
	"github.com/onahvictor/bank/internal/sdk/netutil"
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
	token       auth.Auntenticator
	config      *config.Config
	SessionRepo SessionRepository
}

func NewUserService(ur UserRepository, token auth.Auntenticator, config config.Config, sr SessionRepository) *userService {
	return &userService{
		userRepo:    ur,
		token:       token,
		config:      &config,
		SessionRepo: sr,
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
	SessionID             uuid.UUID
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	User                  *entity.User
}

func (us *userService) Login(ctx context.Context, username, password string, r *http.Request) (*AuthResult, error) {
	user, err := us.userRepo.GetUser(ctx, username)
	fmt.Println(user)

	if err != nil {
		if err == repo.ErrUserNotFound {
			return nil, NewAppError(ErrBadRequest, "user not found", err)
		}
		return nil, NewAppError(ErrInternal, "internal server error:", err)
	}

	if !auth.ComparePassword([]byte(user.HashedPassword), password) {
		return nil, repo.ErrInvalidCredentials
	}

	accessToken, accessPayload, err := us.token.GenerateToken(user.Username, us.config.ACCESS_TOKEN_DURATATION)
	if err != nil {
		return nil, NewAppError(ErrInternal, "token generation failed: %w", err)
	}

	refreshToken, refreshpayload, err := us.token.GenerateToken(user.Username, us.config.REFRESH_TOKEN_DURATION)
	if err != nil {
		return nil, NewAppError(ErrInternal, "token generation failed: %w", err)
	}

	clientIP := netutil.GetClientIP(r)
	session, err := us.SessionRepo.CreateSession(ctx, entity.Session{
		ID:           uuid.MustParse(refreshpayload.ID),
		Username:     user.Username,
		RefreshToken: refreshToken,
		ClientIp:     clientIP,
		UserAgent:    r.UserAgent(),
		IsBlocked:    false,
		ExpiresAt:    refreshpayload.ExpiresAt.Time,
	})
	if err != nil {
		return nil, NewAppError(ErrInternal, "failed to create session", err)
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
		return RenewAccessToken{}, NewAppError(ErrUnauthorized, "session has expired", err)
	}

	session, err := us.SessionRepo.GetSession(ctx, uuid.MustParse(refreshPayload.ID))
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrSessionNotFound):
			return RenewAccessToken{}, NewAppError(ErrNotFound, "session not found", err)
		}
		return RenewAccessToken{}, NewAppError(ErrInternal, "internal server error", err)
	}

	if session.IsBlocked {
		return RenewAccessToken{}, NewAppError(ErrUnauthorized, "session is blocked", err)
	}
	if session.Username != refreshPayload.Username {
		return RenewAccessToken{}, NewAppError(ErrUnauthorized, "incorrect session user", err)
	}

	if session.RefreshToken != refreshToken {
		return RenewAccessToken{}, NewAppError(ErrUnauthorized, "mismatched session token", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return RenewAccessToken{}, NewAppError(ErrUnauthorized, "expired session", err)
	}

	accessToken, accessPayload, err := us.token.GenerateToken(refreshPayload.Username, us.config.ACCESS_TOKEN_DURATATION)
	if err != nil {
		return RenewAccessToken{}, NewAppError(ErrInternal, "internal server error", err)
	}

	return RenewAccessToken{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt.Time,
	}, nil
}
