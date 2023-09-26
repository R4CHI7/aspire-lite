package contract

import (
	"errors"
	"net/http"

	"github.com/r4chi7/aspire-lite/model"
)

type LoanStatusUpdate struct {
	LoanID uint         `json:"loan_id"`
	Status model.Status `json:"status"`
}

func (req *LoanStatusUpdate) Bind(r *http.Request) error {
	if req.LoanID == 0 {
		return errors.New("loan_id is required")
	}

	if req.Status != model.StatusApproved && req.Status != model.StatusPaid && req.Status != model.StatusPending {
		return errors.New("invalid status")
	}

	return nil
}
