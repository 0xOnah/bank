package grpctransport

import (
	"context"

	"github.com/0xOnah/bank/internal/sdk/netutil"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGateWayUserAgent = "grpcgateway-user-agent"
	xForwardedFor        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func extractMetadata(ctx context.Context) *Metadata {
	mdt := &Metadata{}
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		mdt.UserAgent = md[grpcGateWayUserAgent][0]
		ip := md[xForwardedFor][0]
		//prevent storing private ip address as they are no use to us
		if !netutil.IsPrivate(ip) {
			mdt.ClientIP = ip
		}
	}

	peer, ok := peer.FromContext(ctx)
	if ok {
		if !netutil.IsPrivate(peer.Addr.String()) {
			mdt.ClientIP = peer.Addr.String()
		}
	}
	return mdt
}
