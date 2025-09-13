package grpctransport

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryRecoverPanicInterceptor(log *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Str("panic", fmt.Sprintf("%v", r)).
					Bytes("stack", debug.Stack()).
					Msg("panic recoverd")
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

func LoggingInterceptor(log *zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp any, err error) {
		start := time.Now()
		resp, err = handler(ctx, req)

		logger := log.Info()
		statusCode := codes.OK
		if err != nil {
			st, ok := status.FromError(err)
			logger = log.Error().Err(errors.New(st.Message()))
			if ok {
				statusCode = st.Code()
			} else {
				statusCode = codes.Unknown
			}
		}

		defer func() {
			logger.
				Dur("duration_ms", time.Since(start)).
				Str("method", info.FullMethod).
				Str("grpc_type", "unary").
				Int("status_code", int(statusCode)).
				Str("status_text", statusCode.String()).
				Msg("recieved a grpc request")
		}()

		grpc.ChainUnaryInterceptor()
		return resp, err
	}
}
