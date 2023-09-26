package repository

import (
	"context"
	"log"

	"github.com/r4chi7/aspire-lite/model"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func (user User) Create(ctx context.Context, input model.User) (model.User, error) {
	err := user.db.Create(&input).Error
	if err != nil {
		log.Printf("error occurred while saving user in DB: %s", err.Error())
		return model.User{}, err
	}

	return input, nil
}

func NewUser(db *gorm.DB) User {
	return User{db: db}
}
