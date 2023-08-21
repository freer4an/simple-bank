package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mdata := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgent := md.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			mdata.UserAgent = userAgent[0]
		}

		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			mdata.UserAgent = userAgent[0]
		}

		if client := md.Get(xForwardedForHeader); len(client) > 0 {
			mdata.ClientIP = client[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mdata.ClientIP = p.Addr.String()
	}

	return mdata

}
