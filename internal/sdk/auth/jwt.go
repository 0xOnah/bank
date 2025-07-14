package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minKeySize = 32

var (
	ErrTokenGen     = errors.New("failed to generate token")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpired      = errors.New("token has expired")
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(key string) (Auntenticator, error) {
	if len(key) < minKeySize {
		return nil, fmt.Errorf("invalid key size, must be %d characters", minKeySize)
	}
	return &JWTMaker{secretKey: key}, nil
}

func (jt *JWTMaker) GenerateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", ErrTokenGen
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(jt.secretKey))
	if err != nil {
		return "", errors.Join(ErrTokenGen, err)
	}
	return tokenString, nil
}

func (jt *JWTMaker) VerifyToken(token string) (*Payload, error) {
	payload := Payload{}
	parsedToken, err := jwt.ParseWithClaims(token, &payload, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jt.secretKey), nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpired
		}
		return nil, fmt.Errorf("token validation failed")
	}
	if !parsedToken.Valid {
		return nil, ErrInvalidToken
	}
	return &payload, nil
}
