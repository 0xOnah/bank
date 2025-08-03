package grpctransport

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/0xOnah/bank/internal/sdk/validator"
	"github.com/0xOnah/bank/internal/service"
	"github.com/0xOnah/bank/internal/transport/sdk/errorutil"
	"github.com/0xOnah/bank/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (uh *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userValue, err := uh.us.CreateUser(ctx, service.CreateUserInput{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Fullname: req.FullName,
	})

	if err != nil {
		grpcErr := MapValidationErrors(err)
		if grpcErr != nil {
			return nil, grpcErr
		}
		if appErr, ok := err.(*errorutil.AppError); ok {
			slog.Warn("createUser.service", slog.Any("validation failed", appErr.Err.Error()))
			return nil, status.Errorf(errorutil.MapErrorToGRPCStatus(appErr), appErr.Message)
		}
	}
	return &pb.CreateUserResponse{
		User: &pb.User{
			Username: userValue.Username,
			FullName: userValue.FullName,
			Email:    req.Email,
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
			fmt.Println(appErr)
			slog.Warn("loginUser service", slog.Any("app error", appErr.Err.Error()))
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

func MapValidationErrors(err error) error {
	var violations []*errdetails.BadRequest_FieldViolation
	var validatorErr *validator.Validator
	ok := errors.As(err, &validatorErr)
	if ok {
		for key, value := range validatorErr.ErrVal {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       key,
				Description: value,
			})
		}
		badRequest := &errdetails.BadRequest{FieldViolations: violations}
		statusInvalid := status.New(codes.InvalidArgument, "invalid request")
		statusDetails, err := statusInvalid.WithDetails(badRequest)
		if err != nil {
			return statusInvalid.Err()
		}
		return statusDetails.Err()
	}
	return nil
}
