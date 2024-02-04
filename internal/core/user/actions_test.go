package user_test

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go-users-service/internal/core/user"
	repositorymocks "go-users-service/internal/core/user/mocks"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestActions_CreateUser_shouldBeOk(t *testing.T) {
	repo := repositorymocks.NewMockRepository(gomock.NewController(t))
	data := user.Data{
		ID:        uuid.New(),
		FistName:  "test",
		LastName:  "test",
		Birthdate: time.Now(),
		Status:    "ACTIVE",
	}
	repo.EXPECT().SaveUser(gomock.Any(), data).Return(nil)
	actions := user.NewActions(repo)
	err := actions.CreateUser(context.Background(), data)

	assert.Nil(t, err)
}

func TestActions_CreateUser_shouldFail(t *testing.T) {
	repo := repositorymocks.NewMockRepository(gomock.NewController(t))
	data := user.Data{
		ID:        uuid.New(),
		FistName:  "test",
		LastName:  "test",
		Birthdate: time.Now(),
		Status:    "ACTIVE",
	}
	expectedError := errors.New("fail to save")
	repo.EXPECT().SaveUser(gomock.Any(), data).Return(expectedError)
	actions := user.NewActions(repo)
	err := actions.CreateUser(context.Background(), data)

	assert.ErrorContains(t, err, expectedError.Error())
}
