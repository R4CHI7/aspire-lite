package repository

import (
	"context"
	"log"

	"github.com/r4chi7/aspire-lite/model"
	"gorm.io/gorm"
)

type Loan struct {
	db *gorm.DB
}

func (loan Loan) Create(ctx context.Context, input model.Loan) (model.Loan, error) {
	err := loan.db.Create(&input).Error
	if err != nil {
		log.Printf("error occurred while saving loan in DB: %s", err.Error())
		return model.Loan{}, err
	}

	return input, nil
}

func NewLoan(db *gorm.DB) Loan {
	return Loan{db: db}
}
