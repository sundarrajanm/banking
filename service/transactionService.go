package service

import (
	"banking/domain"
	"banking/dto"
	"banking/errs"
	"time"
)

type TransactionService interface {
	NewTransaction(dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError)
}

type DefaultTransactionService struct {
	acctRepo domain.AccountRepository
	txRepo   domain.TransactionRepository
}

func (s DefaultTransactionService) NewTransaction(req dto.NewTransactionRequest) (*dto.NewTransactionResponse, *errs.AppError) {

	err := req.Validate()
	if err != nil {
		return nil, err
	}

	account, err := s.acctRepo.GetAccount(req.AccountId)
	if err != nil {
		return nil, err
	}

	// Do the Transaction - debit or credit
	if req.TransactionType == "withdrawal" {
		if account.Amount < req.Amount {
			return nil, errs.NewValidationError("Withdraw amount is more than the account balance")
		}

		err := s.acctRepo.Debit(req.Amount, req.AccountId)
		if err != nil {
			return nil, err
			// TODO: Work on reverting the debit transaction
		}
	}

	if req.TransactionType == "deposit" {
		err := s.acctRepo.Credit(req.Amount, req.AccountId)
		if err != nil {
			return nil, err
			// TODO: Work on reverting the credit transaction
		}
	}

	// Debit/Credit succeed, record transaction
	tx := domain.Transaction{
		TransactionId:   "",
		AccountId:       req.AccountId,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
		Amount:          req.Amount,
	}

	newTx, err := s.txRepo.Save(tx)
	if err != nil {
		return nil, err
		// TODO: Work on reverting the transaction
	}

	// Get the latest balance
	account, err = s.acctRepo.GetAccount(req.AccountId)
	if err != nil {
		return nil, err
		// TODO: Work on reverting the transaction
	}

	response := dto.NewTransactionResponse{
		NewBalance:    account.Amount,
		TransactionId: newTx.TransactionId,
	}

	return &response, nil
}

func NewTransactionService(acctRepo domain.AccountRepository, txRepo domain.TransactionRepository) DefaultTransactionService {
	return DefaultTransactionService{acctRepo, txRepo}
}
