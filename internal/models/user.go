package models

import (
	"fmt"
	"net/http"
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
