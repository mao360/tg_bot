package payment

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidAmount   = errors.New("invalid amount")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrPaymentFailed   = errors.New("payment failed")
)
