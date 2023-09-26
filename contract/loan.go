package contract

import (
	"errors"
	"net/http"
	"time"
)

type Loan struct {
	Amount float64 `json:"amount"`
	Term   int     `json:"term"`
}

func (loan *Loan) Bind(r *http.Request) error {
	if loan.Amount == 0 {
		return errors.New("amount is required")
	}

	if loan.Term == 0 {
		return errors.New("term is required")
	}

	return nil
}

type LoanResponse struct {
	ID         uint                    `json:"id"`
	Amount     float64                 `json:"amount"`
	Term       int                     `json:"term"`
	Status     string                  `json:"status"`
	CreatedAt  time.Time               `json:"created_at"`
	Repayments []LoanRepaymentResponse `json:"repayments"`
}
