package grpctransport

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/0xOnah/bank/internal/sdk/auth"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationType   = "bearer"
)

func (us *UserHandler) authorization(ctx context.Context) (*auth.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	values := md.Get(strings.ToLower(authorizationHeader))
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	params := strings.Fields(values[0])
	if len(params) != 2 {
		return nil, fmt.Errorf("malformed authorization header format")
	}

	authType := params[0]
	if strings.ToLower(authType) != authorizationType {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	payload, err := us.jwtMaker.VerifyToken(params[1])
	if err != nil {
		if errors.Is(err, auth.ErrExpired) {
			return nil, fmt.Errorf("access token expired")
		}
		return nil, fmt.Errorf("invalid access token")
	}

	return payload, nil
}
