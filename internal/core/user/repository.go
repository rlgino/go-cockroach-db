package user

import (
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	SaveUser(ctx context.Context, user Data) error
	ListUsers(ctx context.Context) ([]Data, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
