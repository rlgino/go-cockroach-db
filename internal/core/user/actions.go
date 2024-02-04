package user

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func NewActions(repo Repository) Actions {
	return &actions{Repo: repo}
}

type Actions interface {
	DeleteUser(ctx context.Context, id uuid.UUID) error
	CreateUser(ctx context.Context, userToCreate Data) error
	FindUser(ctx context.Context, id uuid.UUID) (Data, error)
}

type actions struct {
	Repo Repository
}

func (actions *actions) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := actions.Repo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting user")
	}
	return nil
}

func (actions *actions) CreateUser(ctx context.Context, userToCreate Data) error {
	err := actions.Repo.SaveUser(ctx, userToCreate)
	if err != nil {
		return fmt.Errorf("error creating user %v", err)
	}
	return nil

}

func (actions *actions) FindUser(ctx context.Context, id uuid.UUID) (Data, error) {
	users, err := actions.Repo.FindUserByID(ctx, id)
	if err != nil {
		return Data{}, fmt.Errorf("error find user %s: %v", id.String(), err)
	}
	return users, nil
}
