package service

import (
	"context"
	"errors"
	"testing"

	"github.com/r4chi7/aspire-lite/contract"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

func TestLoanTestSuite(t *testing.T) {
	suite.Run(t, new(LoanTestSuite))
}
