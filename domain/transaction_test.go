package domain

import (
	"banking/dto"
	"net/http"
	"testing"
)

func Test_should_return_error_when_transaction_type_is_not_deposit_or_withdrawal(t *testing.T) {
	request := dto.NewTransactionRequest{
		TransactionType: "Invalid Type",
	}

	err := request.Validate()

	if err.Message != "Transaction type is not withdrawal or deposit" {
		t.Error("Invalid message while testing transaction type")
	}

	if err.Code != http.StatusUnprocessableEntity {
		t.Error("Invalid code while testing transaction type")
	}
}
