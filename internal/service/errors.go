package service

import "errors"

var (
	ErrNumberIsNotCorrect               = errors.New("order number is not correct")
	ErrOrderAlreadyExists               = errors.New("order already exists")
	ErrOrderAlreadyExistsForAnotherUser = errors.New("order already exists for another user")
	ErrInvalidWithdrawSum               = errors.New("invalid withdraw sum")
)
