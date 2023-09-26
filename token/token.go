package token

import (
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

const (
	DefaultExpiry = 24 * time.Hour
)

func Init() {
	tokenAuth = jwtauth.New("HS256", []byte(os.Getenv("TOKEN_SECRET")), nil)
}

func Generate(claims map[string]interface{}) string {
	jwtauth.SetExpiryIn(claims, DefaultExpiry)
	_, token, err := tokenAuth.Encode(claims)
	if err != nil {
		panic(err)
	}

	return token
}

func GetTokenAuth() *jwtauth.JWTAuth {
	return tokenAuth
}
