package controller

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/r4chi7/aspire-lite/contract"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/stretchr/testify/suite"
)

type LoanTestSuite struct {
	suite.Suite
	controller  Loan
	mockService *MockLoanService
}

func (suite *LoanTestSuite) SetupTest() {
	suite.mockService = &MockLoanService{}
	suite.controller = NewLoan(suite.mockService)
}

func (suite *LoanTestSuite) TestCreateHappyFlow() {
	req := httptest.NewRequest(http.MethodPost, "/users/loans", strings.NewReader(`{"amount":10000,"term":2}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), contract.Loan{Amount: 10000, Term: 2}).Return(contract.LoanResponse{
		ID: 1, Amount: 10000, Term: 3, Status: model.StatusPending.String(), Repayments: []contract.LoanRepaymentResponse{
			{
				Amount:  5000,
				DueDate: "2023-10-01",
				Status:  model.StatusPending.String(),
			}, {
				Amount:  5000,
				DueDate: "2023-10-08",
				Status:  model.StatusPending.String(),
			},
		}}, nil)

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusCreated, res.StatusCode)
	suite.Equal(`{"id":1,"amount":10000,"term":3,"status":"PENDING","repayments":[{"amount":5000,"due_date":"2023-10-01","status":"PENDING"},{"amount":5000,"due_date":"2023-10-08","status":"PENDING"}]}
`, string(body))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *LoanTestSuite) TestCreateShouldReturnBadRequestWhenRequestBodyIsIncomplete() {
	req := httptest.NewRequest(http.MethodPost, "/users/loans", strings.NewReader(`{"amount":10000}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusBadRequest, res.StatusCode)
	suite.Equal(`{"status_text":"bad request","message":"term is required"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertNotCalled(suite.T(), "Created")
}

func (suite *LoanTestSuite) TestCreateShouldReturnServerErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodPost, "/users/loans", strings.NewReader(`{"amount":10000,"term":3}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), contract.Loan{Amount: 10000, Term: 3}).Return(contract.LoanResponse{}, errors.New("some error"))

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusInternalServerError, res.StatusCode)
	suite.Equal(`{"status_text":"internal server error","message":"something went wrong, please try again later.."}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertExpectations(suite.T())
}

func TestLoanTestSuite(t *testing.T) {
	suite.Run(t, new(LoanTestSuite))
}
