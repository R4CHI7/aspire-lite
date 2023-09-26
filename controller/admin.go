package controller

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"

	"github.com/r4chi7/aspire-lite/contract"
)

type Admin struct {
	service LoanService
}

func (admin Admin) UpdateLoanStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	isAdmin := claims["is_admin"].(bool)
	if !isAdmin {
		render.Render(w, r, contract.ForbiddenErrorRenderer(errors.New("not admin")))
		return
	}

	input := contract.LoanStatusUpdate{}
	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	err = admin.service.UpdateStatus(ctx, input)
	if err != nil {
		if err == sql.ErrNoRows {
			render.Render(w, r, contract.NotFoundErrorRenderer(errors.New("loan does not exist")))
			return
		}
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)
}

func NewAdmin(service LoanService) Admin {
	return Admin{service: service}
}
