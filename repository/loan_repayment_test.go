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

func (suite *LoanRepaymentTestSuite) TestUpdateHappyFlow() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "loan_repayments" SET "amount"=$1,"status"=$2,"updated_at"=$3 WHERE "id" = $4`)).
		WithArgs(5000, 1, sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.Update(context.Background(), 1, map[string]interface{}{
		"amount": 5000,
		"status": 1,
	})

	suite.Nil(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanRepaymentTestSuite) TestUpdateReturnsErrorWhenDBFails() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "loan_repayments" SET "amount"=$1,"status"=$2,"updated_at"=$3 WHERE "id" = $4`)).
		WithArgs(5000, 1, sqlmock.AnyArg(), 1).WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	err := suite.repo.Update(context.Background(), 1, map[string]interface{}{
		"amount": 5000,
		"status": 1,
	})

	suite.Equal("some error", err.Error())
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanRepaymentTestSuite) TestBulkDeleteHappyFlow() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "loan_repayments" WHERE "loan_repayments"."id" IN ($1,$2,$3)`)).
		WithArgs(1, 2, 3).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.BulkDelete(context.Background(), []uint{1, 2, 3})
	suite.Nil(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanRepaymentTestSuite) TestBulkDeleteReturnsErrorWhenDBFails() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "loan_repayments" WHERE "loan_repayments"."id" IN ($1,$2,$3)`)).
		WithArgs(1, 2, 3).WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	err := suite.repo.BulkDelete(context.Background(), []uint{1, 2, 3})
	suite.Equal("some error", err.Error())
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestLoanRepaymentTestSuite(t *testing.T) {
	suite.Run(t, new(LoanRepaymentTestSuite))
}
