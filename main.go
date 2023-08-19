package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/freer4an/simple-bank/api"
	db "github.com/freer4an/simple-bank/db/sqlc"
	"github.com/freer4an/simple-bank/gapi"
	"github.com/freer4an/simple-bank/pb"
	"github.com/freer4an/simple-bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.InitConfig(".")
	if err != nil {
		log.Fatal("error reading config", err)
	}

	conn, err := sql.Open(config.DB_driver, config.DB_source)
	if err != nil {
		log.Fatal("Connection to db failed:", err)
	}

	store := db.NewStore(conn)
	// go runGinServer(config, store)
	runGrpcServer(config, store)

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("can't create HTTP server: ", err)
	}

	if err := server.Start(config.HttpServerAddr); err != nil {
		log.Fatal("can't start the HTTP server", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("can't create gRPC server: ", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddr)
	if err != nil {
		log.Fatal("cannot create listener: ", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal("can't start the gRPC server: ", err)
	}
}
