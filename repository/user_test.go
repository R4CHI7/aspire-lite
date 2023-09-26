package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserTestSuite struct {
	suite.Suite
	repo User
	mock sqlmock.Sqlmock
}

func (suite *UserTestSuite) SetupTest() {
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

	suite.repo = User{db: db}
	suite.mock = mock
}

func (suite *UserTestSuite) TestCreateHappyFlow() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("email","password","is_admin","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs("test@example.xyz", sqlmock.AnyArg(), false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1))
	suite.mock.ExpectCommit()

	resp, err := suite.repo.Create(context.Background(), model.User{Email: "test@example.xyz", Password: "password"})

	suite.Equal(1, int(resp.ID))
	suite.NoError(err)
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserTestSuite) TestCreateReturnsErrorWhenDBFails() {
	suite.mock.ExpectBegin()
	suite.mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("email","password","is_admin","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs("test@example.xyz", sqlmock.AnyArg(), false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(errors.New("some error"))
	suite.mock.ExpectRollback()

	resp, err := suite.repo.Create(context.Background(), model.User{Email: "test@example.xyz", Password: "password"})

	suite.Empty(resp)
	suite.Error(err, "some error")
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserTestSuite) TestGetByEmailReturnsDataIfExists() {
	// user := model.User{}
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
		WithArgs("test@example.xyz").WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "is_admin"}).
		AddRow(1, "test@example.xyz", "password", false))

	resp, err := suite.repo.GetByEmail(context.Background(), "test@example.xyz")
	suite.NoError(err)
	suite.Equal(1, int(resp.ID))
	suite.Equal("test@example.xyz", resp.Email)
}

func (suite *UserTestSuite) TestGetByEmailReturnsErrorIfDBReturnsError() {
	suite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
		WithArgs("test@example.xyz").WillReturnError(errors.New("some error"))

	resp, err := suite.repo.GetByEmail(context.Background(), "test@example.xyz")
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
