package e2e

import (
	"os"
	"testing"

	"github.com/r4chi7/aspire-lite/database"
	"github.com/r4chi7/aspire-lite/token"
)

func TestMain(m *testing.M) {
	database.InitWithHost("testdb", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	token.Init()
	code := m.Run()
	os.Exit(code)
}
