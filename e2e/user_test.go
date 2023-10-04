package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/r4chi7/aspire-lite/model"
)

func (suite *IntegrationTestSuite) TestUserCreate() {
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"email":"test@example.com","password":"test@123"}`)))
	suite.Nil(err)

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.userController.Create)

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusCreated, rr.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	suite.Nil(err)
	suite.NotEmpty(resp["token"])
	var insertedUser model.User
	suite.db.Where("email = ?", "test@example.com").Find(&insertedUser)
	suite.NotZero(insertedUser.ID)
}

func (suite *IntegrationTestSuite) TestUserLoginHappyFlow() {
	suite.createUser("login@example.com", "test@123", false)
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer([]byte(`{"email":"login@example.com","password":"test@123"}`)))
	suite.Nil(err)

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.userController.Login)

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	suite.Nil(err)
	suite.NotEmpty(resp["token"])
}

func (suite *IntegrationTestSuite) TestUserLoginReturnsErrorWhenUserDoesNotExist() {
	suite.createUser("doesntexist@example.com", "test@123", false)
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer([]byte(`{"email":"login1@example.com","password":"test@123"}`)))
	suite.Nil(err)

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.userController.Login)

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusNotFound, rr.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	suite.Nil(err)
	suite.Equal("user not found", resp["message"])
}

func (suite *IntegrationTestSuite) TestUserLoginReturnsErrorWhenPasswordDoesNotMatch() {
	suite.createUser("password@example.com", "test@123", false)
	req, err := http.NewRequest("POST", "/users/login", bytes.NewBuffer([]byte(`{"email":"password@example.com","password":"test@1235"}`)))
	suite.Nil(err)

	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(suite.userController.Login)

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)

	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	suite.Nil(err)
	suite.Equal("you are unauthorized to perform this action", resp["message"])
}
