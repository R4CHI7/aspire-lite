package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/r4chi7/aspire-lite/contract"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoanTestSuite struct {
	suite.Suite
	controller  Loan
	mockService *MockLoanService
	ctx         context.Context
}

func (suite *LoanTestSuite) SetupTest() {
	suite.mockService = &MockLoanService{}
	suite.controller = NewLoan(suite.mockService)
	token := getToken(map[string]interface{}{
		"user_id": 1.0,
	})
	suite.ctx = jwtauth.NewContext(context.Background(), token, nil)
}

func (suite *LoanTestSuite) TestCreateHappyFlow() {
	req := httptest.NewRequest(http.MethodPost, "/users/loans", strings.NewReader(`{"amount":10000,"term":2}`)).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), uint(1), contract.Loan{Amount: 10000, Term: 2}).Return(contract.LoanResponse{
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
	suite.Equal(`{"id":1,"amount":10000,"term":3,"status":"PENDING","created_at":"0001-01-01T00:00:00Z","repayments":[{"amount":5000,"due_date":"2023-10-01","status":"PENDING"},{"amount":5000,"due_date":"2023-10-08","status":"PENDING"}]}
`, string(body))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *LoanTestSuite) TestCreateShouldReturnBadRequestWhenRequestBodyIsIncomplete() {
	req := httptest.NewRequest(http.MethodPost, "/users/loans", strings.NewReader(`{"amount":10000}`)).WithContext(suite.ctx)
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
	req := httptest.NewRequest(http.MethodPost, "/users/loans", strings.NewReader(`{"amount":10000,"term":3}`)).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), uint(1), contract.Loan{Amount: 10000, Term: 3}).Return(contract.LoanResponse{}, errors.New("some error"))

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

func (suite *LoanTestSuite) TestGetHappyFlow() {
	req := httptest.NewRequest(http.MethodGet, "/users/loans", nil).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("GetByUser", req.Context(), uint(1)).Return([]contract.LoanResponse{{
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
		}}}, nil)

	suite.controller.Get(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusOK, res.StatusCode)
	suite.Equal(`[{"id":1,"amount":10000,"term":3,"status":"PENDING","created_at":"0001-01-01T00:00:00Z","repayments":[{"amount":5000,"due_date":"2023-10-01","status":"PENDING"},{"amount":5000,"due_date":"2023-10-08","status":"PENDING"}]}]
`, string(body))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *LoanTestSuite) TestGetShouldReturnErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodGet, "/users/loans", nil).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("GetByUser", req.Context(), uint(1)).Return([]contract.LoanResponse{}, errors.New("some error"))

	suite.controller.Get(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusInternalServerError, res.StatusCode)
	suite.Equal(`{"status_text":"internal server error","message":"something went wrong, please try again later.."}
`, string(body))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *LoanTestSuite) TestRepayHappyFlow() {
	req := httptest.NewRequest(http.MethodPost, "/users/loans/1/repay", strings.NewReader(`{"amount":5000}`)).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("loanID", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	suite.mockService.On("Repay", mock.Anything, uint(1), uint(1), contract.LoanRepayment{Amount: 5000.0}).Return(nil)

	suite.controller.Repay(w, req)

	res := w.Result()
	suite.Equal(http.StatusOK, res.StatusCode)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *LoanTestSuite) TestRepayReturnsErrorWhenAmountIsZero() {
	req := httptest.NewRequest(http.MethodPost, "/users/loans/1/repay", strings.NewReader(`{"amount":0}`)).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("loanID", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	suite.controller.Repay(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusBadRequest, res.StatusCode)
	suite.Equal(`{"status_text":"bad request","message":"amount should be greater than 0"}
`, string(body))
}

func (suite *LoanTestSuite) TestRepayReturnsErrorWhenServiceReturnsZero() {
	req := httptest.NewRequest(http.MethodPost, "/users/loans/1/repay", strings.NewReader(`{"amount":5000}`)).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("loanID", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	suite.mockService.On("Repay", mock.Anything, uint(1), uint(1), contract.LoanRepayment{Amount: 5000.0}).Return(errors.New("some error"))

	suite.controller.Repay(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusInternalServerError, res.StatusCode)
	suite.Equal(`{"status_text":"internal server error","message":"something went wrong, please try again later.."}
`, string(body))
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
