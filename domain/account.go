package domain

import (
	"banking/dto"
	"banking/errs"
)

type Account struct {
	AccountId   string
	CustomerId  string
	AccountType string
	OpeningDate string
	Amount      float64
	Status      string
}

func (a Account) ToNewAccountResponseDTO() dto.NewAccountResponse {
	return dto.NewAccountResponse{AccountId: a.AccountId}
}

type AccountRepository interface {
	Save(Account) (*Account, *errs.AppError)
}
