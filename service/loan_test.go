package service

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
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

	token := getToken(map[string]interface{}{
		"user_id": 1.0,
	})
	suite.ctx = jwtauth.NewContext(context.Background(), token, nil)
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

	resp, err := suite.service.Create(suite.ctx, input)
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

	resp, err := suite.service.Create(suite.ctx, input)
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

	resp, err := suite.service.Create(suite.ctx, input)
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func getToken(claims map[string]interface{}) jwt.Token {
	tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
	token, _, err := tokenAuth.Encode(claims)
	if err != nil {
		panic(err)
	}
	return token
}

func TestLoanTestSuite(t *testing.T) {
	suite.Run(t, new(LoanTestSuite))
}
