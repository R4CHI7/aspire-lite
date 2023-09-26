package contract

import (
	"errors"
	"net/http"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func (user *User) Bind(r *http.Request) error {
	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}
