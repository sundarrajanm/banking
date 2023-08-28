package dto

type NewTransactionRequest struct {
	AccountId       string
	TransactionType string `json:"transaction_type"`
	Amount          float64
}
