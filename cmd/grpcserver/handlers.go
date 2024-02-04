package main

import (
	"context"
	"github.com/google/uuid"
	"go-users-service/cmd/grpcserver/usersproto"
	"go-users-service/internal/core/user"
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
	id, err := uuid.Parse(request.User)
	if err != nil {
		return nil, err
	}
	user, err := u.actions.FindUser(ctx, id)
	if err != nil {
		return nil, err
	}
	status := usersproto.User_ACTIVE
	if user.Status == "INACTIVE" {
		status = usersproto.User_INACTIVE
	}
	return &usersproto.User{
		Id:       user.ID.String(),
		Name:     user.FistName,
		LastName: user.LastName,
		Datetime: user.Birthdate.Format("2006-01-02"),
		Status:   status,
	}, nil
}
