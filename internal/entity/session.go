package entity

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (ses *Session) IsSessionBlocked() bool {
	return ses.IsBlocked
}

func (ses *Session) UsernameCheck(username string) bool {
	return ses.Username != username
}

func (ses *Session) RefreshTokenCheck(token string) bool {
	return ses.RefreshToken != token
}

func (ses *Session) IsSessionExpired() bool {
	return time.Now().After(ses.ExpiresAt)
}
