package grpctransport

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/0xOnah/bank/internal/db/repo"
	"github.com/0xOnah/bank/internal/sdk/auth"
	"github.com/0xOnah/bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (uh *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := uh.authorization(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	if req.GetUsername() != authPayload.Username {
		return nil, status.Error(codes.PermissionDenied, "you do not have access to this resources")
	}

	//TODO: validate data before
	var password *string
	if req.Password != nil {
		hashed, err := auth.HashPassword(*req.Password)
		if err != nil {
			slog.Error("failed to hash password", slog.Any("error", err))
			return nil, status.Error(codes.Internal, "failed to update user")
		}
		password = &hashed
	}

	if req.Username == "" {
		slog.Error("username not provided", slog.Any("error", err))
		return nil, status.Error(codes.InvalidArgument, "username must be provided")
	}

	user, err := uh.ur.UpdateUser(ctx, repo.UpdateUserParams{
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: password,
		Username:       req.Username,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("user does not exist", slog.Any("error", err))
			return nil, status.Error(codes.NotFound, "user does not exist")
		}
		slog.Error("internal server error", slog.Any("error", err))
		return nil, status.Error(codes.Internal, "internal server error")

	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			Username:          user.Username,
			Email:             user.Email.String(),
			FullName:          user.FullName,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}, nil

}
