package models

import (
	"fmt"
	"net/http"
	"time"
)

type UserCreditials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (uc *UserCreditials) Bind(r *http.Request) error {
	if uc.Login == "" {
		return fmt.Errorf("login field is required")
	}
	if uc.Password == "" {
		return fmt.Errorf("password field is required")
	}
	return nil
}

type Orders struct {
	OrderID    int       `json:"order_id"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
	UserID     int       `json:"user_id"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

type WithdrawBalanceRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
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
	ID          int       `json:"-"`
	UserID      string    `json:"-"`
	OrderNumber string    `json:"order"`
	Amount      int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
