package gapi

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/freer4an/simple-bank/models"
	"github.com/freer4an/simple-bank/pb"
	"github.com/freer4an/simple-bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	accesToken, accesPayload, err := server.tokenMaker.CreateToken(user.Username,
		server.config.AccesTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username,
		server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	md := server.extractMetadata(ctx)
	session := models.Session{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    md.UserAgent,
		ClientIp:     md.ClientIP,
		ExpiresAt:    refreshPayload.ExpiresAt,
	}

	err = server.redis.Set(ctx, req.Username, session, 2*24*time.Hour).Err()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	rsp := &pb.LoginUserResponse{
		Uuid:                session.ID.String(),
		AccesToken:          accesToken,
		AccesTokenExpires:   timestamppb.New(accesPayload.ExpiresAt),
		RefreshToken:        refreshToken,
		RefreshTokenExpires: timestamppb.New(refreshPayload.ExpiresAt),
		User:                convertUser(user),
	}
	return rsp, nil
}
