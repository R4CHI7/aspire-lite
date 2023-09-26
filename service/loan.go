package service

import (
	"context"
	"math"
	"time"

	"github.com/r4chi7/aspire-lite/contract"
	"github.com/r4chi7/aspire-lite/model"
	"gorm.io/datatypes"
)

type Loan struct {
	loanRepository          LoanRepository
	loanRepaymentRepository LoanRepaymentRepository
}

func (loan Loan) Create(ctx context.Context, userID uint, input contract.Loan) (contract.LoanResponse, error) {
	// Create Loan
	loanObj := model.Loan{
		UserID: userID,
		Amount: input.Amount,
		Term:   input.Term,
		Status: model.StatusPending,
	}
	loanObj, err := loan.loanRepository.Create(ctx, loanObj)
	if err != nil {
		return contract.LoanResponse{}, err
	}

	// Create Repayments
	repayments := make([]model.LoanRepayment, 0)
	repaymentAmount := input.Amount / float64(input.Term)
	today := time.Now()
	for i := 1; i <= input.Term; i++ {
		dueDate := today.Add(time.Hour * time.Duration(24*7*i))
		repayments = append(repayments, model.LoanRepayment{
			LoanID:  loanObj.ID,
			Amount:  math.Round(repaymentAmount*100) / 100,
			Status:  model.StatusPending,
			DueDate: datatypes.Date(dueDate),
		})
	}

	err = loan.loanRepaymentRepository.Create(ctx, repayments)
	if err != nil {
		return contract.LoanResponse{}, err
	}

	repaymentResp := make([]contract.LoanRepaymentResponse, 0)
	for _, repayment := range repayments {
		repaymentResp = append(repaymentResp, contract.LoanRepaymentResponse{
			Amount:  repayment.Amount,
			DueDate: time.Time(repayment.DueDate).Format(time.DateOnly),
			Status:  repayment.Status.String(),
		})
	}

	return contract.LoanResponse{
		ID:         loanObj.ID,
		Amount:     loanObj.Amount,
		Term:       loanObj.Term,
		Status:     loanObj.Status.String(),
		CreatedAt:  loanObj.CreatedAt,
		Repayments: repaymentResp,
	}, nil
}

func (loan Loan) GetByUser(ctx context.Context, userID uint) ([]contract.LoanResponse, error) {
	loans, err := loan.loanRepository.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := make([]contract.LoanResponse, 0)
	for _, l := range loans {
		repaymentResp := make([]contract.LoanRepaymentResponse, 0)
		for _, repayment := range l.Repayments {
			repaymentResp = append(repaymentResp, contract.LoanRepaymentResponse{
				Amount:  repayment.Amount,
				DueDate: time.Time(repayment.DueDate).Format(time.DateOnly),
				Status:  repayment.Status.String(),
			})
		}

		resp = append(resp, contract.LoanResponse{
			ID:         l.ID,
			Amount:     l.Amount,
			Term:       l.Term,
			Status:     l.Status.String(),
			CreatedAt:  l.CreatedAt,
			Repayments: repaymentResp,
		})
	}

	return resp, nil
}

func NewLoan(loanRepository LoanRepository, loanRepaymentRepository LoanRepaymentRepository) Loan {
	return Loan{loanRepository: loanRepository, loanRepaymentRepository: loanRepaymentRepository}
}
