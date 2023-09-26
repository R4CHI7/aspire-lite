package controller

import (
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/r4chi7/aspire-lite/contract"
)

type Loan struct {
	service LoanService
}

func (loan Loan) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := contract.Loan{}

	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	resp, err := loan.service.Create(ctx, input)
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
