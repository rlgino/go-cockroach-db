package main

import (
	"context"
	"go-users-service/cmd/grpcserver/usersproto"
	"go-users-service/internal/core/user"
	"log"
)

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usersproto/users.proto

func NewHandlers(repo user.Repository) usersproto.UserServiceServer {
	return &usersServiceServer{
		usersproto.UnimplementedUserServiceServer{},
		user.NewActions(repo),
	}
}

type usersServiceServer struct {
	usersproto.UnimplementedUserServiceServer
	actions user.Actions
}

func (u *usersServiceServer) SearchUser(ctx context.Context, request *usersproto.SearchRequest) (*usersproto.User, error) {
	users, err := u.actions.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("Searching in %d records", len(users))
	for _, data := range users {
		if data.ID.String() == request.User {
			return &usersproto.User{
				Id:     data.ID.String(),
				Name:   data.FistName,
				Status: usersproto.User_ACTIVE,
			}, nil
		}
	}
	log.Printf("FistName with ID %s not found", request.User)
	return &usersproto.User{}, nil
}
