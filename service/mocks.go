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

func (mock *MockUserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	args := mock.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}

type MockLoanRepository struct {
	mock.Mock
}

func (mock *MockLoanRepository) Create(ctx context.Context, loan model.Loan) (model.Loan, error) {
	args := mock.Called(ctx, loan)
	return args.Get(0).(model.Loan), args.Error(1)
}

type MockLoanRepaymentRepository struct {
	mock.Mock
}

func (mock *MockLoanRepaymentRepository) Create(ctx context.Context, repayments []model.LoanRepayment) error {
	args := mock.Called(ctx, repayments)
	return args.Error(0)
}
