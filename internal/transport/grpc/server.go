package grpctransport

import (
	"context"

	"github.com/0xOnah/bank/internal/db/repo"
	"github.com/0xOnah/bank/internal/entity"
	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/internal/service"
	"github.com/0xOnah/bank/pb"
)

type userService interface {
	CreateUser(ctx context.Context, cu service.CreateUserInput) (entity.User, error)
	Login(ctx context.Context, lg service.Logininput) (*service.AuthResult, error)
	RenewAccessToken(ctx context.Context, refreshToken string) (service.RenewAccessToken, error)
}

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	us       userService
	ur       *repo.UserRepo
	jwtMaker auth.Authenticator
}

func NewUserHandler(us userService, ur *repo.UserRepo, jtmaker auth.Authenticator) *UserHandler {
	return &UserHandler{
		us:       us,
		ur:       ur,
		jwtMaker: jtmaker,
	}
}
