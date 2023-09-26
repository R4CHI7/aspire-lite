package controller

import (
	"log"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/r4chi7/aspire-lite/contract"
)

type Loan struct {
	service LoanService
}

func (loan Loan) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := contract.Loan{}

	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	resp, err := loan.service.Create(ctx, uint(claims["user_id"].(float64)), input)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

func NewLoan(service LoanService) Loan {
	return Loan{service: service}
}
