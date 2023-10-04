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

func (suite *IntegrationTestSuite) TestUpdateLoanStatusHappyFlow() {
	userID := suite.createUser("loanadmin@example.com", "password", false)
	userToken := suite.getToken(map[string]interface{}{"user_id": userID})
	loanID := suite.createLoan(userToken, 10000.0, 2)

	adminID := suite.createUser("admin@example.com", "password", true)
	adminToken := suite.getToken(map[string]interface{}{"user_id": adminID, "is_admin": true})
	ctx := jwtauth.NewContext(context.Background(), adminToken, nil)

	input := map[string]interface{}{
		"loan_id": loanID,
		"status":  model.StatusApproved,
	}
	reqBody, _ := json.Marshal(input)
	req, err := http.NewRequest("PATCH", "/admin/loan/status", bytes.NewBuffer(reqBody))
	suite.Nil(err)
	req = req.WithContext(ctx)

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.adminController.UpdateLoanStatus)

	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)

	var loan model.Loan
	suite.db.Find(&loan, loanID)
	suite.Equal(model.StatusApproved, loan.Status)
}

func (suite *IntegrationTestSuite) TestUpdateLoanStatusReturnsForbiddenWhenUserIsNotAdmin() {
	userID := suite.createUser("loanadmin1@example.com", "password", false)
	token := suite.getToken(map[string]interface{}{"user_id": userID, "is_admin": false})
	loanID := suite.createLoan(token, 10000.0, 2)
	ctx := jwtauth.NewContext(context.Background(), token, nil)

	input := map[string]interface{}{
		"loan_id": loanID,
		"status":  model.StatusApproved,
	}
	reqBody, _ := json.Marshal(input)
	req, err := http.NewRequest("PATCH", "/admin/loan/status", bytes.NewBuffer(reqBody))
	suite.Nil(err)
	req = req.WithContext(ctx)

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.adminController.UpdateLoanStatus)

	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusForbidden, rr.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	suite.Nil(err)
	suite.Equal("you are forbidden to perform this action", resp["message"])
}

func (suite *IntegrationTestSuite) TestUpdateLoanStatusReturnsNotFoundWhenLoanDoesNotExist() {
	userID := suite.createUser("loanadmin2@example.com", "password", false)
	userToken := suite.getToken(map[string]interface{}{"user_id": userID})
	suite.createLoan(userToken, 10000.0, 2)

	adminID := suite.createUser("admin1@example.com", "password", true)
	adminToken := suite.getToken(map[string]interface{}{"user_id": adminID, "is_admin": true})
	ctx := jwtauth.NewContext(context.Background(), adminToken, nil)

	input := map[string]interface{}{
		"loan_id": 9999,
		"status":  model.StatusApproved,
	}
	reqBody, _ := json.Marshal(input)
	req, err := http.NewRequest("PATCH", "/admin/loan/status", bytes.NewBuffer(reqBody))
	suite.Nil(err)
	req = req.WithContext(ctx)

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.adminController.UpdateLoanStatus)

	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusNotFound, rr.Code)
	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	suite.Nil(err)
	suite.Equal("loan does not exist", resp["message"])
}
