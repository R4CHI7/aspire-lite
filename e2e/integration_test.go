package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func (suite *IntegrationTestSuite) createUser(email, password string, isAdmin bool) {
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
	suite.db.Where("email = ?", "test@example.com").Find(&insertedUser)
	suite.NotZero(insertedUser.ID)
}
