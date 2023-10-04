package repository

import (
	"context"
	"database/sql"
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

func (loan Loan) GetByUser(ctx context.Context, userID uint) ([]model.Loan, error) {
	var loans []model.Loan
	err := loan.db.Model(&model.Loan{}).Preload("Repayments").Where("user_id = ?", userID).Find(&loans).Error
	if err != nil {
		log.Printf("error occurred while getting loans from DB: %s", err.Error())
		return nil, err
	}

	return loans, nil
}

func (loan Loan) UpdateStatus(ctx context.Context, loanID uint, status model.Status) error {
	loanObj := model.Loan{ID: loanID}
	res := loan.db.Model(&loanObj).Update("status", status)
	if res.Error != nil {
		log.Printf("error occurred while updating loan status: %s", res.Error.Error())
		return res.Error
	}

	if res.RowsAffected == 0 {
		log.Printf("loan not found with ID: %d", loanID)
		return sql.ErrNoRows
	}

	return nil
}

func (loan Loan) GetByID(ctx context.Context, loanID uint) (model.Loan, error) {
	obj := model.Loan{ID: loanID}
	err := loan.db.Model(&model.Loan{}).Preload("Repayments").Find(&obj).Error
	if err != nil {
		log.Printf("error occurred while getting loan by ID: %s", err.Error())
		return model.Loan{}, err
	}

	return obj, nil
}

func NewLoan(db *gorm.DB) Loan {
	return Loan{db: db}
}
