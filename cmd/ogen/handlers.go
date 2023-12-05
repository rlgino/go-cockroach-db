package main

import (
	"context"
	"go-users-service/cmd/ogen/usersvcapi"
	"go-users-service/internal/core/user"
	"net/http"
)

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest -package usersvcapi --target usersvcapi --clean usersvc-oas.yml

func NewHandlers(repo user.Repository) usersvcapi.Handler {
	return userHandlers{
		actions: user.NewActions(repo),
	}
}

type userHandlers struct {
	actions user.Actions
}

func (u userHandlers) DeleteUser(ctx context.Context, params usersvcapi.DeleteUserParams) error {
	err := u.actions.DeleteUser(ctx, params.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (u userHandlers) ListUsers(ctx context.Context) (usersvcapi.Users, error) {
	users, err := u.actions.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	var res = make(usersvcapi.Users, len(users))
	for i, data := range users {
		res[i] = usersvcapi.User{
			ID:   data.ID,
			Name: data.User,
		}
	}
	return res, nil
}

func (u userHandlers) AddUser(ctx context.Context, req *usersvcapi.User) error {
	err := u.actions.CreateUser(ctx, user.Data{
		ID:       req.ID,
		User:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u userHandlers) NewError(_ context.Context, err error) *usersvcapi.ErrorStatusCode {
	return &usersvcapi.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: usersvcapi.Error{
			Message: err.Error(),
		},
	}
}
