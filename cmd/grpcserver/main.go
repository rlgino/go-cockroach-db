package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go-users-service/cmd/grpcserver/usersproto"
	"go-users-service/internal/core/user"
	"go-users-service/internal/persistence"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	conn, cockroachRepository, err := CreateConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	// Start GRPC
	var opts []grpc.ServerOption
	log.Printf("grpc-ping: starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3030"
		log.Printf("Defaulting to port %s", port)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	grpcServer := grpc.NewServer(opts...)
	usersproto.RegisterUserServiceServer(grpcServer, NewHandlers(cockroachRepository))
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func CreateConnection() (*pgx.Conn, user.Repository, error) {
	// Read in connection string
	config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, nil, err
	}
	config.RuntimeParams["application_name"] = "$ docs_simplecrud_gopgx"
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, nil, err
	}
	return conn, persistence.New(conn), nil
}
