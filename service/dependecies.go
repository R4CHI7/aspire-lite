package service

import (
	"context"

	"github.com/r4chi7/aspire-lite/model"
)

type UserRepository interface {
	Create(context.Context, model.User) (model.User, error)
	GetByEmail(context.Context, string) (model.User, error)
}

type LoanRepository interface {
	Create(context.Context, model.Loan) (model.Loan, error)
}

type LoanRepaymentRepository interface {
	Create(context.Context, []model.LoanRepayment) error
}
