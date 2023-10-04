package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/r4chi7/aspire-lite/contract"
)

type Loan struct {
	service LoanService
}

// Create - Creates a new loan
// @Summary This API creates a new loan for the authenticated user
// @Tags loan
// @Accept json
// @Produce json
// @Param event body contract.Loan true "Add loan"
// @Param Authorization header string true "Bearer"
// @Success 200 {object} contract.LoanResponse
// @Router /users/loans [post]
func (loan Loan) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	input := contract.Loan{}

	if err := render.Bind(r, &input); err != nil {
		log.Printf("unable to bind request body: %s", err.Error())
		render.Render(w, r, contract.ErrorRenderer(err))
		return
	}

	resp, err := loan.service.Create(ctx, getUserID(ctx), input)
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, resp)
}

// Create - Gets all loans
// @Summary This API returns all loans for the authenticated user
// @Tags loan
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer"
// @Success 200 {object} []contract.LoanResponse
// @Router /users/loans [get]
func (loan Loan) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := loan.service.GetByUser(ctx, getUserID(ctx))
	if err != nil {
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// Create - Repays a loan
// @Summary This API performs a repayment for the loan
// @Tags loan
// @Accept json
// @Produce json
// @Param event body contract.LoanRepayment true "Add loan"
// @Param Authorization header string true "Bearer"
// @Param loan_id path int true "loan id"
// @Router /users/loans/{loan_id}/repay [post]
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
		if errors.As(err, &contract.RepaymentError{}) {
			render.Render(w, r, contract.ErrorRenderer(err))
			return
		}
		render.Render(w, r, contract.ServerErrorRenderer(err))
		return
	}

	render.Status(r, http.StatusOK)
}

func NewLoan(service LoanService) Loan {
	return Loan{service: service}
}
