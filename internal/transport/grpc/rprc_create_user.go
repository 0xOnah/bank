package grpctransport

import (
	"context"

	"github.com/0xOnah/bank/internal/db/sqlc"
	"github.com/0xOnah/bank/internal/sdk/jobs"
	"github.com/0xOnah/bank/internal/service"
	"github.com/0xOnah/bank/internal/transport/sdk/errorutil"
	"github.com/0xOnah/bank/pb"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (uh *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	var user *sqlc.User
	err := uh.ur.WithTx(ctx, func(q *sqlc.Queries) error {
		var err error
		user, err = q.CreateUser(ctx, sqlc.CreateUserParams{
			Username:       req.Username,
			HashedPassword: req.Password,
			Email:          req.Email,
			FullName:       req.FullName,
		})

		if err != nil {
			grpcErr := MapValidationErrors(err)
			if grpcErr != nil {
				return grpcErr
			}
			if appErr, ok := err.(*errorutil.AppError); ok {
				return status.Errorf(errorutil.MapErrorToGRPCStatus(appErr), appErr.Message)
			}
		}

		//don't do this next time intead use the  cause this could lead to a long lived transaction
		payload := jobs.VerifyEmailPayload{Username: user.Username}

		err = uh.taskqueue.DistributeTaskVerifyEmail(ctx, &payload)
		if err != nil {
			uh.logger.Error().Err(err).Msg("failed to distribute task to send verify email")
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		User: &pb.User{
			Username:          user.Username,
			Email:             user.Email,
			FullName:          user.FullName,
			PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
			CreatedAt:         timestamppb.New(user.CreatedAt),
		},
	}, nil
}

func (uh *UserHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	metadata := extractMetadata(ctx)
	userValue, err := uh.us.Login(ctx, service.Logininput{
		Username:  req.Username,
		Password:  req.Password,
		ClientIP:  metadata.ClientIP,
		UserAgent: metadata.UserAgent,
	})
	if err != nil {
		grpcErr := MapValidationErrors(err)
		if grpcErr != nil {
			return nil, grpcErr
		}
		if appErr, ok := err.(*errorutil.AppError); ok {
			return nil, status.Errorf(errorutil.MapErrorToGRPCStatus(appErr), appErr.Message)
		}
	}
	return &pb.LoginUserResponse{
		SessionId:             userValue.SessionID.String(),
		AccessToken:           userValue.AccessToken,
		AccessTokenExpiresAt:  timestamppb.New(userValue.AccessTokenExpiresAt),
		RefreshToken:          userValue.RefreshToken,
		RefreshTokenExpiresAt: timestamppb.New(userValue.AccessTokenExpiresAt),
		User: &pb.User{
			Username:          userValue.User.Username,
			FullName:          userValue.User.FullName,
			Email:             userValue.User.Email.String(),
			PasswordChangedAt: timestamppb.New(userValue.User.PasswordChangedAt),
			CreatedAt:         timestamppb.New(userValue.User.CreatedAt),
		},
	}, nil
}
