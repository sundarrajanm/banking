package dto

type NewTransactionResponse struct {
	NewBalance    float64 `json:"new_balance"`
	TransactionId string  `json:"transaction_id"`
}
