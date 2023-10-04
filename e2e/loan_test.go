package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/r4chi7/aspire-lite/model"
)

func (suite *IntegrationTestSuite) TestLoanCreateHappyFlow() {
	userID := suite.createUser("loancreate1@example.com", "password", false)
	ctx := jwtauth.NewContext(context.Background(), suite.getToken(map[string]interface{}{"user_id": userID}), nil)

	req, err := http.NewRequest("POST", "/users/loans", bytes.NewBuffer([]byte(`{"amount":10000,"term":2}`)))
	suite.Nil(err)
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.loanController.Create)

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusCreated, rr.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	suite.Nil(err)
	suite.NotEmpty(resp["id"])
	var loan model.Loan
	suite.db.Preload("Repayments").Find(&loan, resp["id"])
	suite.Equal(float64(10000), loan.Amount)
	suite.Equal(2, len(loan.Repayments))
	suite.Equal(float64(5000), loan.Repayments[0].Amount)
}

func (suite *IntegrationTestSuite) TestLoanGetHappyFlow() {
	userID := suite.createUser("loanget@example.com", "password", false)
	token := suite.getToken(map[string]interface{}{"user_id": userID})
	ctx := jwtauth.NewContext(context.Background(), token, nil)

	loan1ID := suite.createLoan(token, 10000.0, 2)
	loan2ID := suite.createLoan(token, 15000.0, 3)

	req, err := http.NewRequest("GET", "/users/loans", nil)
	suite.Nil(err)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.loanController.Get)

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	var resp []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	suite.Nil(err)
	suite.Equal(loan1ID, resp[0]["id"])
	suite.Equal(loan2ID, resp[1]["id"])
}

func (suite *IntegrationTestSuite) TestRepayHappyFlowEqualRepayments() {
	userID := suite.createUser("loanrepay@example.com", "password", false)
	token := suite.getToken(map[string]interface{}{"user_id": userID})
	loanID := suite.createLoan(token, 10000.0, 2)
	ctx := jwtauth.NewContext(context.Background(), token, nil)

	suite.approveLoan(loanID)

	// First repayment
	req, err := http.NewRequest("POST", fmt.Sprintf("/users/loans/%d/repay", uint(loanID)), bytes.NewBuffer([]byte(`{"amount":5000}`)))
	suite.Nil(err)
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("loanID", fmt.Sprintf("%d", int(loanID)))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.loanController.Repay)

	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)

	var loan model.Loan
	suite.db.Preload("Repayments").Find(&loan, loanID)
	suite.Equal(model.StatusPaid, loan.Repayments[1].Status)

	// Second repayment
	req, err = http.NewRequest("POST", fmt.Sprintf("/users/loans/%d/repay", uint(loanID)), bytes.NewBuffer([]byte(`{"amount":5000}`)))
	suite.Nil(err)
	req = req.WithContext(ctx)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	req.Header.Add("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)

	suite.db.Preload("Repayments").Find(&loan, loanID)
	suite.Equal(model.StatusPaid, loan.Status)
	suite.Equal(model.StatusPaid, loan.Repayments[0].Status)
	suite.Equal(model.StatusPaid, loan.Repayments[1].Status)
}

func (suite *IntegrationTestSuite) TestRepayHappyFlowUnEqualRepayments() {
	userID := suite.createUser("loanrepay1@example.com", "password", false)
	token := suite.getToken(map[string]interface{}{"user_id": userID})
	loanID := suite.createLoan(token, 10000.0, 2)
	ctx := jwtauth.NewContext(context.Background(), token, nil)

	suite.approveLoan(loanID)

	// First repayment
	req, err := http.NewRequest("POST", fmt.Sprintf("/users/loans/%d/repay", uint(loanID)), bytes.NewBuffer([]byte(`{"amount":6000}`)))
	suite.Nil(err)
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("loanID", fmt.Sprintf("%d", int(loanID)))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.loanController.Repay)

	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)

	// Assert that the pending repayment's amount got updated
	var loan model.Loan
	suite.db.Preload("Repayments").Find(&loan, loanID)
	for _, repay := range loan.Repayments {
		switch repay.Status {
		case model.StatusPending:
			suite.Equal(4000.0, repay.Amount)
		case model.StatusPaid:
			suite.Equal(6000.0, repay.Amount)
		}
	}

	// Second repayment
	req, err = http.NewRequest("POST", fmt.Sprintf("/users/loans/%d/repay", uint(loanID)), bytes.NewBuffer([]byte(`{"amount":4000}`)))
	suite.Nil(err)
	req = req.WithContext(ctx)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	req.Header.Add("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)

	suite.db.Preload("Repayments").Find(&loan, loanID)
	suite.Equal(model.StatusPaid, loan.Status)
	suite.Equal(model.StatusPaid, loan.Repayments[0].Status)
	suite.Equal(model.StatusPaid, loan.Repayments[1].Status)
}

func (suite *IntegrationTestSuite) TestRepayDeletesPendingPaymentIfUserPaysOffLoanPrematurely() {
	userID := suite.createUser("loanrepay2@example.com", "password", false)
	token := suite.getToken(map[string]interface{}{"user_id": userID})
	loanID := suite.createLoan(token, 10000.0, 3)
	ctx := jwtauth.NewContext(context.Background(), token, nil)

	suite.approveLoan(loanID)

	// First repayment
	req, err := http.NewRequest("POST", fmt.Sprintf("/users/loans/%d/repay", uint(loanID)), bytes.NewBuffer([]byte(`{"amount":4000}`)))
	suite.Nil(err)
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("loanID", fmt.Sprintf("%d", int(loanID)))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.loanController.Repay)

	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)

	// Assert that the pending repayment's amount got updated
	var loan model.Loan
	suite.db.Preload("Repayments").Find(&loan, loanID)
	suite.Equal(3, len(loan.Repayments))
	for _, repay := range loan.Repayments {
		switch repay.Status {
		case model.StatusPending:
			suite.Equal(3000.0, repay.Amount)
		case model.StatusPaid:
			suite.Equal(4000.0, repay.Amount)
		}
	}

	// Second repayment
	req, err = http.NewRequest("POST", fmt.Sprintf("/users/loans/%d/repay", uint(loanID)), bytes.NewBuffer([]byte(`{"amount":6000}`)))
	suite.Nil(err)
	req = req.WithContext(ctx)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	req.Header.Add("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)

	// Assert that the pending repayment's amount got deleted
	suite.db.Preload("Repayments").Find(&loan, loanID)
	suite.Equal(2, len(loan.Repayments))
	suite.Equal(model.StatusPaid, loan.Status)
}
