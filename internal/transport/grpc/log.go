package grpctransport

import (
	"context"

	"google.golang.org/grpc"
)

func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {

	return nil, nil
}
