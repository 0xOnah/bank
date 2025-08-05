package grpctransport

import (
	"context"
	"errors"
	"net"

	"github.com/0xOnah/bank/internal/sdk/netutil"
	"github.com/0xOnah/bank/internal/sdk/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	grpcGateWayUserAgent = "grpcgateway-user-agent"
	xForwardedFor        = "x-forwarded-for"
	userAgent            = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func extractMetadata(ctx context.Context) *Metadata {
	mdt := &Metadata{}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if ua := md[grpcGateWayUserAgent]; len(ua) > 0 {
			mdt.UserAgent = ua[0]
		}

		if ua := md[userAgent]; len(ua) > 0 {
			mdt.UserAgent = ua[0]
		}

		
		ui := md[xForwardedFor]
		if len(ui) > 0 {
			//prevent storing private ip address as they are no use to us
			if !netutil.IsPrivate(ui[0]) {
				mdt.ClientIP = ui[0]
			}
		}
	}

	if mdt.ClientIP == "" {
		peer, ok := peer.FromContext(ctx)
		if ok {
			host, _, err := net.SplitHostPort(peer.Addr.String())
			if err == nil && !netutil.IsPrivate(host) {
				mdt.ClientIP = peer.Addr.String()
			}
		}
	}
	return mdt
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
