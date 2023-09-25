package main

import (
	"os"

	"github.com/r4chi7/aspire-lite/database"
)

func main() {
	database.Init(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
}
