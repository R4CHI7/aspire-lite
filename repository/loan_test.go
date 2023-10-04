package repository

import (
	"context"
	"database/sql"
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

type LoanTestSuite struct {
	suite.Suite
	repo Loan
	mock sqlmock.Sqlmock
}

func (suite *LoanTestSuite) SetupTest() {
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

	suite.repo = Loan{db: db}
	suite.mock = mock
}

func (suite *LoanTestSuite) TestCreateHappyFlow() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "loans" ("user_id","amount","term","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(1, 10000.0, 2, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()

	resp, err := suite.repo.Create(context.Background(), model.Loan{UserID: 1, Amount: 10000, Term: 2, Status: model.StatusPending})

	suite.Equal(1, int(resp.ID))
	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanTestSuite) TestCreateReturnsErrorWhenDBFails() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "loans" ("user_id","amount","term","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
		WithArgs(1, 10000.0, 2, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	resp, err := suite.repo.Create(context.Background(), model.Loan{UserID: 1, Amount: 10000, Term: 2, Status: model.StatusPending})

	suite.Empty(resp)
	suite.Error(err, "some error")
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanTestSuite) TestGetByUserHappyFlow() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "loans" WHERE user_id = $1`)).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "amount", "term", "status", "created_at"}).AddRow(1, 1, 10000, 2, 0, time.Now()))

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "loan_repayments" WHERE "loan_repayments"."loan_id" = $1`)).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "loan_id", "amount", "due_date", "status", "created_at"}).
				AddRow(1, 1, 10000, datatypes.Date(time.Now()), 0, time.Now()).
				AddRow(1, 1, 10000, datatypes.Date(time.Now()), 0, time.Now()))

	resp, err := suite.repo.GetByUser(context.Background(), 1)
	suite.NoError(err)
	suite.Equal(1, len(resp))
	suite.Equal(2, len(resp[0].Repayments))
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanTestSuite) TestGetByUserShouldReturnErrorIfDBFails() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "loans" WHERE user_id = $1`)).
		WithArgs(1).
		WillReturnError(errors.New("some error"))

	resp, err := suite.repo.GetByUser(context.Background(), 1)
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanTestSuite) TestUpdateStatusHappyFlow() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "loans" SET "status"=$1,"updated_at"=$2 WHERE "id" = $3`)).
		WithArgs(1, sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 1))
	suite.mock.ExpectCommit()

	err := suite.repo.UpdateStatus(context.Background(), 1, 1)
	suite.Nil(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanTestSuite) TestUpdateStatusReturnsNoRowsErrorIfNoRowsAffected() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "loans" SET "status"=$1,"updated_at"=$2 WHERE "id" = $3`)).
		WithArgs(1, sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.mock.ExpectCommit()

	err := suite.repo.UpdateStatus(context.Background(), 1, 1)
	suite.Equal(sql.ErrNoRows, err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanTestSuite) TestGetByIDHappyFlow() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "loans" WHERE "loans"."id" = $1`)).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "amount", "term", "status", "created_at"}).AddRow(1, 1, 10000, 2, 0, time.Now()))

	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "loan_repayments" WHERE "loan_repayments"."loan_id" = $1`)).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "loan_id", "amount", "due_date", "status", "created_at"}).
				AddRow(1, 1, 10000, datatypes.Date(time.Now()), 0, time.Now()).
				AddRow(1, 1, 10000, datatypes.Date(time.Now()), 0, time.Now()))

	resp, err := suite.repo.GetByID(context.Background(), 1)
	suite.Nil(err)
	suite.Equal(uint(1), resp.ID)
	suite.Equal(2, len(resp.Repayments))
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *LoanTestSuite) TestGetByIDShouldReturnErrorIfDBFails() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "loans" WHERE "loans"."id" = $1`)).
		WithArgs(1).
		WillReturnError(errors.New("some error"))

	resp, err := suite.repo.GetByID(context.Background(), 1)
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestLoanTestSuite(t *testing.T) {
	suite.Run(t, new(LoanTestSuite))
}
