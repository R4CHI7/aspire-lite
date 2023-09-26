package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/r4chi7/aspire-lite/contract"
	"github.com/stretchr/testify/suite"
)

type AdminTestSuite struct {
	suite.Suite
	controller  Admin
	mockService *MockLoanService
	ctx         context.Context
}

func (suite *AdminTestSuite) SetupTest() {
	suite.mockService = &MockLoanService{}
	suite.controller = NewAdmin(suite.mockService)
	token := getToken(map[string]interface{}{
		"user_id":  1.0,
		"is_admin": true,
	})
	suite.ctx = jwtauth.NewContext(context.Background(), token, nil)
}

func (suite *AdminTestSuite) TestUpdateLoanStatusHappyFlow() {
	req := httptest.NewRequest(http.MethodPatch, "/admin/loan/status", strings.NewReader(`{"loan_id":1,"status":1}`)).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockService.On("UpdateStatus", suite.ctx, contract.LoanStatusUpdate{LoanID: 1, Status: 1}).Return(nil)

	suite.controller.UpdateLoanStatus(w, req)
	res := w.Result()

	suite.Equal(http.StatusOK, res.StatusCode)
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *AdminTestSuite) TestUpdateLoanStatusReturnsForbiddenWhenUserIsNotAdmin() {
	token := getToken(map[string]interface{}{
		"user_id":  1.0,
		"is_admin": false,
	})
	suite.ctx = jwtauth.NewContext(context.Background(), token, nil)
	req := httptest.NewRequest(http.MethodPatch, "/admin/loan/status", strings.NewReader(`{"loan_id":1,"status":1}`)).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.controller.UpdateLoanStatus(w, req)
	res := w.Result()

	suite.Equal(http.StatusForbidden, res.StatusCode)
}

func (suite *AdminTestSuite) TestUpdateLoanStatusReturnsErrorWhenRequestBodyIsIncomplete() {
	req := httptest.NewRequest(http.MethodPatch, "/admin/loan/status", strings.NewReader(`{"status":1}`)).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.controller.UpdateLoanStatus(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}

	suite.Equal(http.StatusBadRequest, res.StatusCode)
	suite.Equal(`{"status_text":"bad request","message":"loan_id is required"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertNotCalled(suite.T(), "UpdateStatus")
}

func (suite *AdminTestSuite) TestUpdateLoanStatusReturnsErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodPatch, "/admin/loan/status", strings.NewReader(`{"loan_id":1,"status":1}`)).WithContext(suite.ctx)
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.mockService.On("UpdateStatus", suite.ctx, contract.LoanStatusUpdate{LoanID: 1, Status: 1}).Return(errors.New("some error"))

	suite.controller.UpdateLoanStatus(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}

	suite.Equal(http.StatusBadRequest, res.StatusCode)
	suite.Equal(`{"status_text":"bad request","message":"some error"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertExpectations(suite.T())
}

func TestAdminTestSuite(t *testing.T) {
	suite.Run(t, new(AdminTestSuite))
}
