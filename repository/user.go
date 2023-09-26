package repository

import (
	"context"
	"database/sql"
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

func (user User) GetByEmail(ctx context.Context, email string) (model.User, error) {
	u := model.User{}
	res := user.db.Where("email = ?", email).Find(&u)
	if res.Error != nil {
		log.Printf("error occurred while getting user availability from DB: %s", res.Error.Error())
		return model.User{}, res.Error
	}

	if res.RowsAffected == 0 {
		log.Printf("user not found for email: %s", email)
		return model.User{}, sql.ErrNoRows
	}

	return u, nil
}

func NewUser(db *gorm.DB) User {
	return User{db: db}
}
