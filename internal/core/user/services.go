package user

import (
	"context"
	"errors"
	"go-cockroach/internal/persistence"
)

type Deleter struct {
	repo *persistence.CockroachRepository
}

func (deleter *Deleter) deleteUser(ctx context.Context, userToDelete Data) error {
	err := deleter.repo.DeleteUser(ctx, userToDelete.ID)
	if err != nil {
		return errors.New("error deleting user")
	}
	return nil
}

type Creator struct {
	repo *persistence.CockroachRepository
}

func (handler *Creator) createUser(ctx context.Context, userToCreate Data) error {
	err := handler.repo.SaveUser(ctx, userToCreate)
	if err != nil {
		return errors.New("error creating user")
	}
	return nil

}

type Lister struct {
	repo *persistence.CockroachRepository
}

func (handler *Lister) listUsers(ctx context.Context) ([]Data, error) {
	users, err := handler.repo.ListUsers(ctx)
	if err != nil {
		return nil, errors.New("error listing user")
	}
	return users, nil
}
