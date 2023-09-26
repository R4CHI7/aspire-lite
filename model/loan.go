package model

import "time"

type status int

const (
	StatusPending  status = 0
	StatusApproved status = 1
	StatusPaid     status = 2
)

func (s status) String() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusApproved:
		return "APPROVED"
	case StatusPaid:
		return "DELETED"
	}
	return ""
}

type Loan struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	Term      int       `gorm:"not null"`
	Status    status    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Repayments []LoanRepayment
}
