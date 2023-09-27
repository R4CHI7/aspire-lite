package controller

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
)

func getUserID(ctx context.Context) uint {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0
	}

	return uint(claims["user_id"].(float64))
}
