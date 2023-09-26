package model

import (
	"time"

	"gorm.io/datatypes"
)

type LoanRepayment struct {
	ID        uint           `gorm:"primaryKey"`
	LoanID    uint           `gorm:"not null"`
	Amount    float64        `gorm:"not null"`
	DueDate   datatypes.Date `gorm:"not null"`
	Status    status         `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
}
