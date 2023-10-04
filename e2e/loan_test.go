package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

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
