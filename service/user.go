package service

import (
	"context"

	"github.com/r4chi7/aspire-lite/contract"
	"github.com/r4chi7/aspire-lite/model"
	"github.com/r4chi7/aspire-lite/token"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	repository UserRepository
}

func (user User) Create(ctx context.Context, input contract.User) (contract.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 8)
	if err != nil {
		return contract.UserResponse{}, err
	}
	userObj := model.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		IsAdmin:  input.IsAdmin,
	}

	userObj, err = user.repository.Create(ctx, userObj)
	if err != nil {
		return contract.UserResponse{}, err
	}

	authToken := token.Generate(map[string]interface{}{
		"user_id":  userObj.ID,
		"is_admin": userObj.IsAdmin,
	})

	return contract.UserResponse{ID: userObj.ID, Token: authToken}, nil
}

func (user User) Login(ctx context.Context, input contract.UserLogin) (contract.UserResponse, error) {
	userObj, err := user.repository.GetByEmail(ctx, input.Email)
	if err != nil {
		return contract.UserResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userObj.Password), []byte(input.Password))
	if err != nil {
		return contract.UserResponse{}, err
	}

	authToken := token.Generate(map[string]interface{}{
		"user_id":  userObj.ID,
		"is_admin": userObj.IsAdmin,
	})

	return contract.UserResponse{ID: userObj.ID, Token: authToken}, nil

}

func NewUser(repo UserRepository) User {
	return User{repository: repo}
}
