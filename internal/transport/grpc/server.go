package grpctransport

import (
	"context"

	"github.com/0xOnah/bank/internal/entity"
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
	us userService
}

func NewUserHandler(us userService) *UserHandler {
	return &UserHandler{
		us: us,
	}
}
