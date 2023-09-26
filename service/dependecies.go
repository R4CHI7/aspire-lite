package service

import (
	"context"

	"github.com/r4chi7/aspire-lite/model"
)

type UserRepository interface {
	Create(context.Context, model.User) (model.User, error)
}
