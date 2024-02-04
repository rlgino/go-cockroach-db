package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go-users-service/cmd/ogen/usersvcapi"
	"go-users-service/internal/core/logger"
	"go-users-service/internal/core/user"
	"net/http"
	"time"
)

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest -package usersvcapi --target usersvcapi --clean usersvc-oas.yml

func NewHandlers(repo user.Repository, logger logger.Logger) usersvcapi.Handler {
	return userHandlers{
		logger:  logger,
		actions: user.NewActions(repo),
	}
}

type userHandlers struct {
	logger  logger.Logger
	actions user.Actions
}

var fields = map[string]interface{}{
	"cluster-uuid": uuid.New().String(),
}

func (u userHandlers) DeleteUser(ctx context.Context, params usersvcapi.DeleteUserParams) error {
	u.logger.Info(fmt.Sprintf("Deleting user %s", params.UserID), fields)
	err := u.actions.DeleteUser(ctx, params.UserID)
	if err != nil {
		u.logger.Error(fmt.Sprintf("Error deleting user %s: %v", params.UserID, err), fields)
		return err
	}
	return nil
}

func (u userHandlers) FindUser(ctx context.Context, params usersvcapi.FindUserParams) (*usersvcapi.User, error) {
	u.logger.Info("Find user "+params.UserID.String(), fields)
	data, err := u.actions.FindUser(ctx, params.UserID)
	if err != nil {
		u.logger.Error(fmt.Sprintf("Error Listing users: %v", err), fields)
		return nil, err
	}
	return &usersvcapi.User{
		ID:        data.ID,
		Name:      data.FistName,
		LastName:  data.LastName,
		Birthdate: data.Birthdate.Format("2006-01-02"),
		Status:    usersvcapi.UserStatus(data.Status),
	}, nil
}

func (u userHandlers) AddUser(ctx context.Context, req *usersvcapi.User) error {
	u.logger.Info(fmt.Sprintf("Creating user %v", req), fields)
	birthdate, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		return err
	}
	err = req.Status.Validate()
	if err != nil {
		return err
	}
	err = u.actions.CreateUser(ctx, user.Data{
		ID:        req.ID,
		FistName:  req.Name,
		LastName:  req.LastName,
		Birthdate: birthdate,
		Status:    string(req.Status),
	})
	if err != nil {
		u.logger.Error(fmt.Sprintf("Error creating user %v: %v", req, err), fields)
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
