package controller

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/r4chi7/aspire-lite/contract"
	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	controller  User
	mockService *MockUserService
}

func (suite *UserTestSuite) SetupTest() {
	suite.mockService = &MockUserService{}
	suite.controller = NewUser(suite.mockService)

}

func (suite *UserTestSuite) TestCreateHappyFlow() {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"email":"test@example.xyz","password":"password"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), contract.User{Email: "test@example.xyz", Password: "password"}).Return(contract.UserResponse{ID: 1, Token: "token"}, nil)

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusCreated, res.StatusCode)
	suite.Equal(`{"id":1,"token":"token"}
`, string(body))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *UserTestSuite) TestCreateShouldReturnBadRequestWhenRequestBodyIsIncomplete() {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"email":"test@example.xyz"}`))
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
	suite.Equal(`{"status_text":"bad request","message":"password is required"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertNotCalled(suite.T(), "Created")
}

func (suite *UserTestSuite) TestCreateShouldReturnServerErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"email":"test@example.xyz","password":"password"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Create", req.Context(), contract.User{Email: "test@example.xyz", Password: "password"}).Return(contract.UserResponse{}, errors.New("some error"))

	suite.controller.Create(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusInternalServerError, res.StatusCode)
	suite.Equal(`{"status_text":"internal server error","message":"some error"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *UserTestSuite) TestLoginHappyFlow() {
	req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email":"test@example.xyz","password":"password"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Login", req.Context(), contract.UserLogin{Email: "test@example.xyz", Password: "password"}).Return(contract.UserResponse{ID: 1, Token: "token"}, nil)

	suite.controller.Login(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusOK, res.StatusCode)
	suite.Equal(`{"id":1,"token":"token"}
`, string(body))
	suite.mockService.AssertExpectations(suite.T())
}

func (suite *UserTestSuite) TestLoginShouldReturnBadRequestWhenRequestBodyIsIncomplete() {
	req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email":"test@example.xyz"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.controller.Login(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusBadRequest, res.StatusCode)
	suite.Equal(`{"status_text":"bad request","message":"password is required"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertNotCalled(suite.T(), "Created")
}

func (suite *UserTestSuite) TestLoginShouldReturnServerErrorWhenServiceReturnsError() {
	req := httptest.NewRequest(http.MethodPost, "/users/login", strings.NewReader(`{"email":"test@example.xyz","password":"password"}`))
	req.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.mockService.On("Login", req.Context(), contract.UserLogin{Email: "test@example.xyz", Password: "password"}).Return(contract.UserResponse{}, errors.New("some error"))

	suite.controller.Login(w, req)

	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		suite.Error(errors.New("expected error to be nil got"), err)
	}
	suite.Equal(http.StatusInternalServerError, res.StatusCode)
	suite.Equal(`{"status_text":"internal server error","message":"some error"}
`, string(body)) // This newline is needed because chi returns the response ending with a \n
	suite.mockService.AssertExpectations(suite.T())
}

func TestUserTest(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
