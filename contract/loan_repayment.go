package contract

import (
	"errors"
	"net/http"
)

type LoanRepayment struct {
	Amount float64 `json:"amount"`
}

func (req LoanRepayment) Bind(r *http.Request) error {
	if req.Amount <= 0.0 {
		return errors.New("amount should be greater than 0")
	}

	return nil
}

type LoanRepaymentResponse struct {
	Amount  float64 `json:"amount"`
	DueDate string  `json:"due_date"`
	Status  string  `json:"status"`
}
