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

func (mock *MockUserService) Login(ctx context.Context, input contract.UserLogin) (contract.UserResponse, error) {
	args := mock.Called(ctx, input)
	return args.Get(0).(contract.UserResponse), args.Error(1)
}

type MockLoanService struct {
	mock.Mock
}

func (mock *MockLoanService) Create(ctx context.Context, userID uint, input contract.Loan) (contract.LoanResponse, error) {
	args := mock.Called(ctx, userID, input)
	return args.Get(0).(contract.LoanResponse), args.Error(1)
}

func (mock *MockLoanService) GetByUser(ctx context.Context, userID uint) ([]contract.LoanResponse, error) {
	args := mock.Called(ctx, userID)
	return args.Get(0).([]contract.LoanResponse), args.Error(1)
}

func (mock *MockLoanService) UpdateStatus(ctx context.Context, input contract.LoanStatusUpdate) error {
	args := mock.Called(ctx, input)
	return args.Error(0)
}
