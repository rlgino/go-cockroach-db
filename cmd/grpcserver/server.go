package grpcserver

import (
	"fmt"
	"go-users-service/cmd/grpcServer/usersproto"
	"go-users-service/internal/core/user"
	"google.golang.org/grpc"
	"log"
	"net"
)

func New(repo user.Repository) *Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	usersproto.RegisterUserServiceServer(grpcServer, NewHandlers(repo))
	return &Server{
		grpcServer,
	}
}

type Server struct {
	grpcServer *grpc.Server
}

func (srv Server) ListenAndServe(port string) {
	log.Println("Listening GRPC in ", port)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}
	if err := srv.grpcServer.Serve(lis); err != nil {
		log.Println(err)
	}
}
