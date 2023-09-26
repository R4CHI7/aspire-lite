package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/r4chi7/aspire-lite/contract"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/r4chi7/aspire-lite/token"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserTestSuite struct {
	suite.Suite
	service            User
	mockUserRepository *MockUserRepository
	ctx                context.Context
}

func (suite *UserTestSuite) SetupSuite() {
	os.Setenv("TOKEN_SECRET", "secret")
	token.Init()
}

func (suite *UserTestSuite) SetupTest() {
	suite.mockUserRepository = &MockUserRepository{}
	suite.service = NewUser(suite.mockUserRepository)
	suite.ctx = context.Background()
}

func (suite *UserTestSuite) TestCreateHappyFlow() {
	input := contract.User{
		Email:    "test@example.xyz",
		Password: "password",
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 8)
	expectedResp := model.User{
		ID:        1,
		Email:     "test@example.xyz",
		Password:  string(hash),
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mockUserRepository.On("Create", suite.ctx, mock.Anything).Return(expectedResp, nil)

	resp, err := suite.service.Create(suite.ctx, input)
	suite.Nil(err)
	suite.Equal(uint(1), resp.ID)
}

func (suite *UserTestSuite) TestCreateShouldReturnErrorIfRepositoryFails() {
	input := contract.User{
		Email:    "test@example.xyz",
		Password: "password",
	}
	suite.mockUserRepository.On("Create", suite.ctx, mock.Anything).Return(model.User{}, errors.New("some error"))

	resp, err := suite.service.Create(suite.ctx, input)
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func (suite *UserTestSuite) TestLoginHappyFlow() {
	input := contract.UserLogin{
		Email:    "test@example.xyz",
		Password: "password",
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 8)
	expectedResp := model.User{
		ID:        1,
		Email:     "test@example.xyz",
		Password:  string(hash),
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mockUserRepository.On("GetByEmail", suite.ctx, "test@example.xyz").Return(expectedResp, nil)

	resp, err := suite.service.Login(suite.ctx, input)
	suite.Nil(err)
	suite.Equal(uint(1), resp.ID)
}

func (suite *UserTestSuite) TestLoginShouldReturnErrorIfRepositoryFails() {
	input := contract.UserLogin{
		Email:    "test@example.xyz",
		Password: "password",
	}
	suite.mockUserRepository.On("GetByEmail", suite.ctx, "test@example.xyz").Return(model.User{}, errors.New("some error"))

	resp, err := suite.service.Login(suite.ctx, input)
	suite.Equal("some error", err.Error())
	suite.Empty(resp)
}

func (suite *UserTestSuite) TestLoginShouldReturnErrorIfPasswordsDontMatch() {
	input := contract.UserLogin{
		Email:    "test@example.xyz",
		Password: "password",
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte("password1"), 8)
	expectedResp := model.User{
		ID:        1,
		Email:     "test@example.xyz",
		Password:  string(hash),
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	suite.mockUserRepository.On("GetByEmail", suite.ctx, "test@example.xyz").Return(expectedResp, nil)

	resp, err := suite.service.Login(suite.ctx, input)
	suite.Equal(bcrypt.ErrMismatchedHashAndPassword, err)
	suite.Empty(resp)
}

func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
