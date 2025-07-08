package auth

import (
	"time"
)

type Auntenticator interface {
	GenerateToken(name string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
