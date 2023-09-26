package main

import (
	"net/http"
	"os"

	"github.com/r4chi7/aspire-lite/database"
	"github.com/r4chi7/aspire-lite/server"
	"github.com/r4chi7/aspire-lite/token"
)

func main() {
	database.Init(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	token.Init()

	r := server.Init()
	http.ListenAndServe(":8080", r)
}
