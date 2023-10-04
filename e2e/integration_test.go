package e2e

import (
	"testing"

	"github.com/r4chi7/aspire-lite/controller"
	"github.com/r4chi7/aspire-lite/database"
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
