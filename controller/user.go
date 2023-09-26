package controller

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/r4chi7/aspire-lite/contract"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	service UserService
}

func (user User) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := contract.User{}

	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	resp, err := user.service.Create(ctx, input)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

func (user User) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := contract.UserLogin{}

	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	resp, err := user.service.Login(ctx, input)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			render.Render(w, r, contract.NotFoundErrorRenderer(errors.New("user not found")))
			return
		case bcrypt.ErrMismatchedHashAndPassword:
			render.Render(w, r, contract.UnauthorizedErrorRenderer(err))
			return
		}
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func NewUser(service UserService) User {
	return User{service: service}
}
