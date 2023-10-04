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
	GetByUser(context.Context, uint) ([]model.Loan, error)
	UpdateStatus(context.Context, uint, model.Status) error
	GetByID(context.Context, uint) (model.Loan, error)
}

type LoanRepaymentRepository interface {
	Create(context.Context, []model.LoanRepayment) error
	Update(context.Context, uint, map[string]interface{}) error
	BulkDelete(context.Context, []uint) error
}
