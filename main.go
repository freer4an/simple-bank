package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"net"
	"net/http"

	db "github.com/freer4an/simple-bank/db/sqlc"
	"github.com/freer4an/simple-bank/gapi"
	"github.com/freer4an/simple-bank/pb"
	"github.com/freer4an/simple-bank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.InitConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("error reading config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DB_driver, config.DB_source)
	if err != nil {
		log.Fatal().Err(err).Msg("Connection to db failed")
	}

	runDBmigration(config.MigrattionURL, config.DB_source)

	store := db.NewStore(conn)

	go runGatewayServer(config, store)
	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create gRPC server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal().Err(err).Msg("can't start the gRPC server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create gRPC server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server); err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	fs := http.FileServer(http.Dir("./docs/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

	listener, err := net.Listen("tcp", config.HttpServerAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener: ")
	}

	log.Info().Msgf("start gRPC-gateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	if err = http.Serve(listener, handler); err != nil {
		log.Fatal().Err(err).Msg("can't start the gRCP-gateway server")
	}
}

// func runGinServer(config util.Config, store db.Store) {
// 	server, err := api.NewServer(config, store)
// 	if err != nil {
// 		log.Fatal().Err(err).Msg("can't create HTTP server: ", err)
// 	}

// 	if err := server.Start(config.HttpServerAddr); err != nil {
// 		log.Fatal().Err(err).Msg("can't start the HTTP server", err)
// 	}
// }

func runDBmigration(migrationUrl, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("migration error")
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("migration up error")
	}

	log.Info().Msg("migration succes")
}
