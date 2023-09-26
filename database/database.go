package database

import (
	"fmt"

	"github.com/r4chi7/aspire-lite/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	HOST = "database"
	PORT = 5432
)

var db *gorm.DB

func Init(username, password, database string) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Loan{})
	if err != nil {
		panic(err)
	}
}

func Get() *gorm.DB {
	return db
}
