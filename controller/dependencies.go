package controller

import (
	"context"

	"github.com/r4chi7/aspire-lite/contract"
)

type UserService interface {
	Create(context.Context, contract.User) (contract.UserResponse, error)
}
