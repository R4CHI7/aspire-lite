package service

import (
	"context"

	"github.com/r4chi7/aspire-lite/model"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (mock *MockUserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	args := mock.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}
