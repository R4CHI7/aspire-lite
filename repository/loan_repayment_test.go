package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LoanRepaymentTestSuite struct {
	suite.Suite
	repo LoanRepayment
	mock sqlmock.Sqlmock
}

func (suite *LoanRepaymentTestSuite) SetupTest() {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		suite.NoError(err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	suite.NoError(err)

	suite.repo = LoanRepayment{db: db}
	suite.mock = mock
}

func (suite *LoanRepaymentTestSuite) TestCreateHappyFlow() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "loan_repayments" ("loan_id","amount","due_date","status","created_at","updated_at")
	VALUES ($1,$2,$3,$4,$5,$6),($7,$8,$9,$10,$11,$12) RETURNING "id"`)).
		WithArgs(1, 5000.0, sqlmock.AnyArg(), 0, sqlmock.AnyArg(), sqlmock.AnyArg(),
			1, 5000.0, sqlmock.AnyArg(), 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()

	err := suite.repo.Create(context.Background(), []model.LoanRepayment{
		{
			LoanID:  1,
			Amount:  5000,
			Status:  model.StatusPending,
			DueDate: datatypes.Date(time.Now().Add(24 * time.Hour)),
		},
		{
			LoanID:  1,
			Amount:  5000,
			Status:  model.StatusPending,
			DueDate: datatypes.Date(time.Now().Add(48 * time.Hour)),
		},
	})

	suite.Nil(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanRepaymentTestSuite) TestCreateReturnsErrorWhenDBFails() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "loan_repayments" ("loan_id","amount","due_date","status","created_at","updated_at")
	VALUES ($1,$2,$3,$4,$5,$6),($7,$8,$9,$10,$11,$12) RETURNING "id"`)).
		WithArgs(1, 5000.0, sqlmock.AnyArg(), 0, sqlmock.AnyArg(), sqlmock.AnyArg(),
			1, 5000.0, sqlmock.AnyArg(), 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	err := suite.repo.Create(context.Background(), []model.LoanRepayment{
		{
			LoanID:  1,
			Amount:  5000,
			Status:  model.StatusPending,
			DueDate: datatypes.Date(time.Now().Add(24 * time.Hour)),
		},
		{
			LoanID:  1,
			Amount:  5000,
			Status:  model.StatusPending,
			DueDate: datatypes.Date(time.Now().Add(48 * time.Hour)),
		},
	})

	suite.Equal("some error", err.Error())
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestLoanRepaymentTestSuite(t *testing.T) {
	suite.Run(t, new(LoanRepaymentTestSuite))
}
