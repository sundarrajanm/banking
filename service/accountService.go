package service

import (
	"banking/domain"
	"banking/dto"
	"banking/errs"
	"time"
)

type AccountService interface {
	NewAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
	NewTransaction(dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repo domain.AccountRepository
}

func (s DefaultAccountService) NewAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	a := domain.Account{
		AccountId:   "",
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format("2006-01-02 15:04:05"),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}

	newAcc, err := s.repo.Save(a)
	if err != nil {
		return nil, err
	}

	response := newAcc.ToNewAccountResponseDTO()
	return &response, nil
}

func (s DefaultAccountService) NewTransaction(req dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError) {

	err := req.Validate()
	if err != nil {
		return nil, err
	}

	tx := domain.Transaction{
		TransactionId:   "",
		AccountId:       req.AccountId,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
		Amount:          req.Amount,
	}

	if tx.IsWithdrawal() {
		account, err := s.repo.GetAccount(req.AccountId)
		if err != nil {
			return nil, err
		}

		if !account.CanWithdraw(req.Amount) {
			return nil, errs.NewValidationError("Withdraw amount is more than the account balance")
		}
	}

	newTx, err := s.repo.DoTransaction(tx)
	if err != nil {
		return nil, err
	}

	response := newTx.ToResponseDTO()

	return &response, nil
}

func NewAccountService(repo domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{repo}
}
