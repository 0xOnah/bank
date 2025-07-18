package auth

import (
	"time"
)

type Auntenticator interface {
	GenerateToken(name string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
