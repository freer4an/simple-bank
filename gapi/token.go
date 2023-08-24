package gapi

import (
	"context"
	"fmt"

	"github.com/freer4an/simple-bank/models"
	"github.com/freer4an/simple-bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RenewAcces(ctx context.Context, req *pb.RefreshTokenReq) (*pb.RefreshTokenRsp, error) {

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	var session models.Session
	err = server.redis.Get(ctx, refreshPayload.Username).Scan(&session)
	if err != nil {
		// if errors.Is(err, redis.) {
		// 	ctx.JSON(http.StatusNotFound, errorResponse(err))
		// 	return
		// }
		return nil, status.Errorf(codes.Internal, "token: %v", err)

	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		return nil, status.Error(codes.Unauthenticated, err.Error())

	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	accesToken, accesPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username,
		server.config.AccesTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())

	}

	rsp := &pb.RefreshTokenRsp{
		AccessToken:        accesToken,
		AccessTokenExpires: timestamppb.New(accesPayload.ExpiresAt),
	}

	return rsp, nil
}
