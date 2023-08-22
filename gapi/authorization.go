package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/freer4an/simple-bank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authHeaderKey  = "authorization"
	authTypeBearer = "bearer"
	authPayloadKey = "authorization_payload"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authHeaderKey)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := strings.Fields(values[0])
	if len(authHeader) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	if authHeader[0] != authTypeBearer {
		return nil, fmt.Errorf("unsupported authorization header type")
	}

	accesToken := authHeader[1]
	payload, err := server.tokenMaker.VerifyToken(accesToken)
	if err != nil {
		return nil, fmt.Errorf("invalid acces token: %s", err)
	}

	return payload, nil
}
