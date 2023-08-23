package gapi

import (
	"fmt"

	db "github.com/freer4an/simple-bank/db/sqlc"
	"github.com/freer4an/simple-bank/pb"
	"github.com/freer4an/simple-bank/token"
	"github.com/freer4an/simple-bank/util"
	"github.com/redis/go-redis/v9"
)

// banking service server
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	redis      *redis.Client
}

// NewServer creates new gRPC server
func NewServer(config util.Config, store db.Store, redis *redis.Client) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("error token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		redis:      redis,
	}

	return server, nil
}

// func (server *Server) Start(addres string) error {
// 	return server.router.Run(addres)
// }

// func errorResponse(err error) gin.H {
// 	return gin.H{"error": err.Error()}
// }
