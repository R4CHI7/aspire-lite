package service

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/r4chi7/aspire-lite/contract"
	"github.com/r4chi7/aspire-lite/model"
	"gorm.io/datatypes"
)

type Loan struct {
	loanRepository          LoanRepository
	loanRepaymentRepository LoanRepaymentRepository
}

func (loan Loan) Create(ctx context.Context, input contract.Loan) (contract.LoanResponse, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		log.Printf("could not get claims from context: %s", err.Error())
		return contract.LoanResponse{}, err
	}

	// Create Loan
	loanObj := model.Loan{
		UserID: uint(claims["user_id"].(float64)),
		Amount: input.Amount,
		Term:   input.Term,
		Status: model.StatusPending,
	}
	loanObj, err = loan.loanRepository.Create(ctx, loanObj)
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

	return contract.LoanResponse{
		ID:     loanObj.ID,
		Amount: loanObj.Amount,
		Term:   loanObj.Term,
	}, nil
}

func NewLoan(loanRepository LoanRepository, loanRepaymentRepository LoanRepaymentRepository) Loan {
	return Loan{loanRepository: loanRepository, loanRepaymentRepository: loanRepaymentRepository}
}
