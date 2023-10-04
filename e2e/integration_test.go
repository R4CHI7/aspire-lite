package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/r4chi7/aspire-lite/controller"
	"github.com/r4chi7/aspire-lite/database"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/r4chi7/aspire-lite/repository"
	"github.com/r4chi7/aspire-lite/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	db              *gorm.DB
	userController  controller.User
	adminController controller.Admin
	loanController  controller.Loan
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.db = database.Get()
	userRepository := repository.NewUser(suite.db)
	loanRepository := repository.NewLoan(suite.db)
	loanRepaymentRepository := repository.NewLoanRepayment(suite.db)

	suite.userController = controller.NewUser(service.NewUser(userRepository))
	suite.loanController = controller.NewLoan(service.NewLoan(loanRepository, loanRepaymentRepository))
	suite.adminController = controller.NewAdmin(service.NewLoan(loanRepository, loanRepaymentRepository))
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

// Some helper methods

func (suite *IntegrationTestSuite) createUser(email, password string, isAdmin bool) float64 {
	input := map[string]interface{}{
		"email":    email,
		"password": password,
		"is_admin": isAdmin,
	}
	reqBody, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(reqBody)))
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
	suite.db.Where("email = ?", email).Find(&insertedUser)
	suite.NotZero(insertedUser.ID)

	return resp["id"].(float64)
}

func (suite *IntegrationTestSuite) createLoan(token jwt.Token, amount float64, term int) float64 {
	ctx := jwtauth.NewContext(context.Background(), token, nil)
	input := map[string]interface{}{
		"amount": amount,
		"term":   term,
	}
	reqBody, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/users/loans", bytes.NewBuffer(reqBody))
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

	return resp["id"].(float64)
}

func (suite *IntegrationTestSuite) approveLoan(id float64) {
	loanObj := model.Loan{ID: uint(id)}
	res := suite.db.Model(&loanObj).Update("status", model.StatusApproved)

	suite.Equal(1, int(res.RowsAffected))
	suite.Nil(res.Error)
}

func (suite *IntegrationTestSuite) getToken(claims map[string]interface{}) jwt.Token {
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("TOKEN_SECRET")), nil)
	token, _, err := tokenAuth.Encode(claims)
	if err != nil {
		panic(err)
	}
	return token
}
