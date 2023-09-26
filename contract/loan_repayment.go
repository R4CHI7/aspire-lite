package contract

type LoanRepaymentResponse struct {
	Amount  float64 `json:"amount"`
	DueDate string  `json:"due_date"`
	Status  string  `json:"status"`
}
