package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/r4chi7/aspire-lite/contract"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
)

type LoanTestSuite struct {
	suite.Suite
	service                     Loan
	mockLoanRepository          *MockLoanRepository
	mockLoanRepaymentRepository *MockLoanRepaymentRepository
	ctx                         context.Context
}

func (suite *LoanTestSuite) SetupTest() {
	suite.mockLoanRepository = &MockLoanRepository{}
	suite.mockLoanRepaymentRepository = &MockLoanRepaymentRepository{}
	suite.service = NewLoan(suite.mockLoanRepository, suite.mockLoanRepaymentRepository)
	suite.ctx = context.Background()
}

func (suite *LoanTestSuite) TestCreateHappyFlow() {
	input := contract.Loan{
		Amount: 10000,
		Term:   2,
	}

	expectedLoanObj := model.Loan{
		ID:     1,
		UserID: 1,
		Amount: 10000,
		Term:   2,
		Status: model.StatusPending,
	}

	suite.mockLoanRepository.On("Create", suite.ctx, model.Loan{
		UserID: 1,
		Amount: 10000,
		Term:   2,
		Status: model.StatusPending,
	}).Return(expectedLoanObj, nil)

	suite.mockLoanRepaymentRepository.On("Create", suite.ctx, mock.Anything).Return(nil)

	resp, err := suite.service.Create(suite.ctx, uint(1), input)
	suite.Nil(err)
	suite.Equal(uint(1), resp.ID)
}

func (suite *LoanTestSuite) TestCreateShouldReturnErrorWhenLoanRepositoryFails() {
	input := contract.Loan{
		Amount: 10000,
		Term:   2,
	}

	expectedLoanObj := model.Loan{
		ID:     1,
		UserID: 1,
		Amount: 10000,
		Term:   2,
		Status: model.StatusPending,
	}

	suite.mockLoanRepository.On("Create", suite.ctx, model.Loan{
		UserID: 1,
		Amount: 10000,
		Term:   2,
		Status: model.StatusPending,
	}).Return(expectedLoanObj, nil)

	suite.mockLoanRepaymentRepository.On("Create", suite.ctx, mock.Anything).Return(errors.New("some error"))

	resp, err := suite.service.Create(suite.ctx, uint(1), input)
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func (suite *LoanTestSuite) TestCreateShouldReturnErrorWhenLoanRepaymentRepositoryFails() {
	input := contract.Loan{
		Amount: 10000,
		Term:   2,
	}

	suite.mockLoanRepository.On("Create", suite.ctx, model.Loan{
		UserID: 1,
		Amount: 10000,
		Term:   2,
		Status: model.StatusPending,
	}).Return(model.Loan{}, errors.New("some error"))

	resp, err := suite.service.Create(suite.ctx, uint(1), input)
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func (suite *LoanTestSuite) TestGetByUserHappyFlow() {
	now := time.Now()
	suite.mockLoanRepository.On("GetByUser", suite.ctx, uint(1)).Return([]model.Loan{
		{
			ID:        1,
			UserID:    1,
			Amount:    10000,
			Term:      2,
			Status:    model.StatusPending,
			CreatedAt: now,
			Repayments: []model.LoanRepayment{
				{
					ID:      1,
					LoanID:  1,
					Amount:  5000,
					DueDate: datatypes.Date(now),
					Status:  model.StatusPending,
				},
				{
					ID:      2,
					LoanID:  1,
					Amount:  5000,
					DueDate: datatypes.Date(now),
					Status:  model.StatusPending,
				},
			},
		}, {
			ID:        2,
			UserID:    1,
			Amount:    12000,
			Term:      2,
			Status:    model.StatusApproved,
			CreatedAt: now,
			Repayments: []model.LoanRepayment{
				{
					ID:      3,
					LoanID:  2,
					Amount:  6000,
					DueDate: datatypes.Date(now),
					Status:  model.StatusPending,
				},
				{
					ID:      4,
					LoanID:  2,
					Amount:  6000,
					DueDate: datatypes.Date(now),
					Status:  model.StatusPending,
				},
			},
		},
	}, nil)

	resp, err := suite.service.GetByUser(suite.ctx, uint(1))
	suite.Nil(err)
	suite.Equal(2, len(resp))
}

func (suite *LoanTestSuite) TestGetByUserShouldReturnErrorIfRepositoryFails() {
	suite.mockLoanRepository.On("GetByUser", suite.ctx, uint(1)).Return([]model.Loan{}, errors.New("some error"))

	resp, err := suite.service.GetByUser(suite.ctx, uint(1))
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func (suite *LoanTestSuite) TestUpdateStatusHappyFlow() {
	suite.mockLoanRepository.On("UpdateStatus", suite.ctx, uint(1), model.StatusApproved).Return(nil)

	err := suite.service.UpdateStatus(suite.ctx, contract.LoanStatusUpdate{LoanID: uint(1), Status: model.StatusApproved})
	suite.Nil(err)
	suite.mockLoanRepository.AssertExpectations(suite.T())
}

func (suite *LoanTestSuite) TestUpdateStatusReturnsErrorWhenRepositoryReturnsError() {
	suite.mockLoanRepository.On("UpdateStatus", suite.ctx, uint(1), model.StatusApproved).Return(errors.New("some error"))

	err := suite.service.UpdateStatus(suite.ctx, contract.LoanStatusUpdate{LoanID: uint(1), Status: model.StatusApproved})
	suite.Equal("some error", err.Error())
	suite.mockLoanRepository.AssertExpectations(suite.T())
}

func TestLoanTestSuite(t *testing.T) {
	suite.Run(t, new(LoanTestSuite))
}
