package main

import (
	"net/http"
	"os"

	"github.com/r4chi7/aspire-lite/database"
	"github.com/r4chi7/aspire-lite/server"
)

func main() {
	database.Init(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	r := server.Init()
	http.ListenAndServe(":8080", r)
}
