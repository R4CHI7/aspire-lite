package database

import (
	"fmt"

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

	err = db.AutoMigrate()
	if err != nil {
		panic(err)
	}
}

func Get() *gorm.DB {
	return db
}
