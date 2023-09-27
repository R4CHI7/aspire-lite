package controller

import (
	"context"

	"github.com/r4chi7/aspire-lite/contract"
)

type UserService interface {
	Create(context.Context, contract.User) (contract.UserResponse, error)
	Login(context.Context, contract.UserLogin) (contract.UserResponse, error)
}

type LoanService interface {
	Create(context.Context, uint, contract.Loan) (contract.LoanResponse, error)
	GetByUser(context.Context, uint) ([]contract.LoanResponse, error)
	UpdateStatus(context.Context, contract.LoanStatusUpdate) error
	Repay(context.Context, uint, uint, contract.LoanRepayment) error
}
