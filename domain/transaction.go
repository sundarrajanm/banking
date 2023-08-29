package domain

import "banking/dto"

type Transaction struct {
	TransactionId   string `db:"transaction_id"`
	AccountId       string `db:"account_id"`
	Amount          float64
	TransactionType string `db:"transaction_type"`
	TransactionDate string `db:"transaction_date"`
}

func (t Transaction) IsWithdrawal() bool {
	return t.TransactionType == "withdrawal"
}

func (t Transaction) ToResponseDTO() dto.NewTransactionResponse {
	return dto.NewTransactionResponse{
		NewBalance:    t.Amount,
		TransactionId: t.TransactionId,
	}
}
