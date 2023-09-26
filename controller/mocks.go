package controller

import (
	"context"

	"github.com/r4chi7/aspire-lite/contract"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (mock *MockUserService) Create(ctx context.Context, input contract.User) (contract.UserResponse, error) {
	args := mock.Called(ctx, input)
	return args.Get(0).(contract.UserResponse), args.Error(1)
}
