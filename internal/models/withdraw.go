package models

import (
	"fmt"
	"net/http"
	"time"
)

type WithdrawBalanceRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

func (uc *WithdrawBalanceRequest) Bind(r *http.Request) error {
	if uc.Order == "" {
		return fmt.Errorf("login field is required")
	}
	if uc.Sum == 0 {
		return fmt.Errorf("password field is required")
	}
	return nil
}

type WithdrawBalance struct {
	WithdrawalID int       `json:"withdrawal_id"`
	UserID       string    `json:"user_id"`
	OrderNumber  string    `json:"order"`
	Amount       float32   `json:"sum"`
	ProcessedAt  time.Time `json:"processed_at"`
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
