package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/onahvictor/bank/internal/db/sqlc"
	"github.com/onahvictor/bank/internal/entity"
)

var (
	ErrFailedToCreateSession = errors.New("failed to create session")
	ErrSessionNotFound       = errors.New("session not found")
)

type sessionRepo struct {
	db *sqlc.SQLStore
}

func toEntitySession(s *sqlc.Session) *entity.Session {
	return &entity.Session{
		ID:           s.ID,
		Username:     s.Username,
		RefreshToken: s.RefreshToken,
		UserAgent:    s.UserAgent,
		ClientIp:     s.ClientIp,
		IsBlocked:    s.IsBlocked,
		ExpiresAt:    s.ExpiresAt,
	}

}

func NewSessionRepo(db *sqlc.SQLStore) *sessionRepo {
	return &sessionRepo{db: db}
}

func (s *sessionRepo) CreateSession(ctx context.Context, arg entity.Session) (*entity.Session, error) {
	result, err := s.db.CreateSession(ctx, sqlc.CreateSessionParams{
		ID:           arg.ID,
		Username:     arg.Username,
		RefreshToken: arg.RefreshToken,
		UserAgent:    arg.UserAgent,
		ClientIp:     arg.ClientIp,
		IsBlocked:    arg.IsBlocked,
		ExpiresAt:    arg.ExpiresAt,
	})
	return toEntitySession(result), err
}

func (s *sessionRepo) GetSession(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	session, err := s.db.GetSession(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSessionNotFound
		}
	}
	return toEntitySession(session), nil
}
