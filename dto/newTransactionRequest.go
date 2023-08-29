package dto

import (
	"banking/errs"
	"strings"
)

type NewTransactionRequest struct {
	AccountId       string
	TransactionType string `json:"transaction_type"`
	Amount          float64
}

func (r NewTransactionRequest) Validate() *errs.AppError {

	if strings.ToLower(r.TransactionType) != "withdrawal" && strings.ToLower(r.TransactionType) != "deposit" {
		return errs.NewValidationError("Transaction type is not withdrawal or deposit")
	}

	if r.Amount <= 0 {
		return errs.NewValidationError("Transaction amount should be greater than 0")
	}

	return nil
}
