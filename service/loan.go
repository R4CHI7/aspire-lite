package service

import (
	"context"
	"fmt"
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

func (loan Loan) UpdateStatus(ctx context.Context, input contract.LoanStatusUpdate) error {
	return loan.loanRepository.UpdateStatus(ctx, input.LoanID, input.Status)
}

func (loan Loan) Repay(ctx context.Context, userID, loanID uint, input contract.LoanRepayment) error {
	loanObj, err := loan.loanRepository.GetByID(ctx, loanID)
	if err != nil {
		return err
	}

	// Check loan belongs to the user
	if loanObj.UserID != userID {
		return contract.NewRepaymentError("loan does not belong to the user")
	}

	// Loan's status should be APPROVED
	if loanObj.Status != model.StatusApproved {
		return contract.NewRepaymentError("loan is not approved")
	}

	pendingRepayments := make([]model.LoanRepayment, 0)
	var paidAmount float64
	for _, repay := range loanObj.Repayments {
		switch repay.Status {
		case model.StatusPending:
			pendingRepayments = append(pendingRepayments, repay)
		case model.StatusPaid:
			paidAmount += repay.Amount
		}
	}

	repayment := pendingRepayments[0]
	// If amount is less than repayment's amount, return error
	if input.Amount < repayment.Amount {
		return contract.NewRepaymentError(fmt.Sprintf("amount should be at least %.2f", repayment.Amount))
	}

	// If this is the last repayment
	if len(pendingRepayments) == 1 {
		if input.Amount > repayment.Amount {
			return contract.NewRepaymentError(fmt.Sprintf("amount should be %.2f", repayment.Amount))
		}

		// Update repayment's status
		err = loan.loanRepaymentRepository.Update(ctx, repayment.ID, map[string]interface{}{
			"amount": input.Amount,
			"status": model.StatusPaid,
		})
		if err != nil {
			return err
		}

		// Update loan's status
		err = loan.loanRepository.UpdateStatus(ctx, loanID, model.StatusPaid)
		if err != nil {
			return err
		}
		return nil
	}

	// If this was not the last repayment
	// Check if user is trying to pay more than required amount.
	if paidAmount+input.Amount > loanObj.Amount {
		return contract.NewRepaymentError(fmt.Sprintf("amount should be less than %.2f", (loanObj.Amount - paidAmount)))
	}

	// Update this repayment's amount and status
	repaymentAmount := repayment.Amount
	err = loan.loanRepaymentRepository.Update(ctx, repayment.ID, map[string]interface{}{
		"amount": input.Amount,
		"status": model.StatusPaid,
	})
	if err != nil {
		return err
	}

	// If user gave amount greather than required, we need to update handle pending repayments
	if input.Amount > repaymentAmount {
		paidAmount += input.Amount
		pendingRepayments = pendingRepayments[1:]

		// If user has paid the complete loan's amount in this repayment, mark loan as paid and delete rest of the repayments
		if paidAmount == loanObj.Amount {
			repaymentsToBeDeleted := make([]uint, 0)
			for _, repayment := range pendingRepayments {
				repaymentsToBeDeleted = append(repaymentsToBeDeleted, repayment.ID)
			}
			err = loan.loanRepaymentRepository.BulkDelete(ctx, repaymentsToBeDeleted)
			if err != nil {
				return err
			}

			err = loan.loanRepository.UpdateStatus(ctx, loanID, model.StatusPaid)
			if err != nil {
				return err
			}

			return nil
		}

		// Update amount for pending repayments.
		pendingAmount := (loanObj.Amount - paidAmount) / float64(len(pendingRepayments))
		for _, r := range pendingRepayments {
			err = loan.loanRepaymentRepository.Update(ctx, r.ID, map[string]interface{}{
				"amount": math.Round(pendingAmount*100) / 100,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewLoan(loanRepository LoanRepository, loanRepaymentRepository LoanRepaymentRepository) Loan {
	return Loan{loanRepository: loanRepository, loanRepaymentRepository: loanRepaymentRepository}
}
