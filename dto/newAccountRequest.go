package dto

import (
	"banking/errs"
	"strings"
)

type NewAccountRequest struct {
	CustomerId  string  `json:"customer_id"`
	AccountType string  `json:"account_type"`
	Amount      float64 `json:"amount"`
}

func (r NewAccountRequest) Validate() *errs.AppError {

	if r.Amount < 5000 {
		return errs.NewValidationError("Amount should be greater than 5000 to create new account")
	}

	if strings.ToLower(r.AccountType) != "checking" && strings.ToLower(r.AccountType) != "saving" {
		return errs.NewValidationError("Only checking or saving account types are supported")
	}

	return nil
}
