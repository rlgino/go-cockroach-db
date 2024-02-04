package user

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockgen -source repository.go -destination mocks/repository.go -package repositorymocks Repository

type Repository interface {
	SaveUser(ctx context.Context, user Data) error
	FindUserByID(ctx context.Context, id uuid.UUID) (Data, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
