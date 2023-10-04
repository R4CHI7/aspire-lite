package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/r4chi7/aspire-lite/model"
)

func (suite *IntegrationTestSuite) TestUser() {
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
