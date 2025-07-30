package auth

import (
	"time"
)

type Authenticator interface {
	GenerateToken(name string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
