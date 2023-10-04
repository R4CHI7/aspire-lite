package model

import "time"

type Status int

const (
	StatusPending  Status = 0
	StatusApproved Status = 1
	StatusPaid     Status = 2
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusApproved:
		return "APPROVED"
	case StatusPaid:
		return "PAID"
	}
	return ""
}

type Loan struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	Term      int       `gorm:"not null"`
	Status    Status    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Repayments []LoanRepayment
}
