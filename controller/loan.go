package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

func (loan Loan) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	resp, err := loan.service.GetByUser(ctx, uint(claims["user_id"].(float64)))
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (loan Loan) Repay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := getUserID(ctx)
	if userID == 0 {
		render.Render(w, r, contract.ServerErrorRenderer(errors.New("invalid token")))
		return
	}

	loanIDs := chi.URLParam(r, "loanID")
	if loanIDs == "" {
		render.Render(w, r, contract.ErrorRenderer(errors.New("loan ID is required")))
		return
	}
	loanID, err := strconv.Atoi(loanIDs)
	if err != nil {
		render.Render(w, r, contract.ErrorRenderer(errors.New("invalid loan ID")))
	}

	input := contract.LoanRepayment{}
	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	err = loan.service.Repay(ctx, userID, uint(loanID), input)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)
}

func NewLoan(service LoanService) Loan {
	return Loan{service: service}
}
