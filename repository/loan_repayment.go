package repository

import (
	"context"
	"log"

	"github.com/r4chi7/aspire-lite/model"
	"gorm.io/gorm"
)

type LoanRepayment struct {
	db *gorm.DB
}

func (repayment LoanRepayment) Create(ctx context.Context, repayments []model.LoanRepayment) error {
	err := repayment.db.Create(repayments).Error
	if err != nil {
		log.Printf("error occurred while inserting repayments: %s", err.Error())
		return err
	}
	return nil
}

func NewLoanRepayment(db *gorm.DB) LoanRepayment {
	return LoanRepayment{db: db}
}
