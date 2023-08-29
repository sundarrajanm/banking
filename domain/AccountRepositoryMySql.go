package domain

import (
	"banking/errs"
	"banking/logger"
	"database/sql"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type AccountRepositoryMySql struct {
	db *sqlx.DB
}

func (a AccountRepositoryMySql) GetAccount(accountId string) (*Account, *errs.AppError) {
	GetAccountSql := "select * from accounts where account_id = ?"
	var account Account
	err := a.db.Get(&account, GetAccountSql, accountId)
	if err == sql.ErrNoRows {
		return nil, errs.NewNotFoundError("Account not found")
	} else if err != nil {
		logger.Error("Unexpected database error while fetching account: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	return &account, nil
}

func (a AccountRepositoryMySql) Save(account Account) (*Account, *errs.AppError) {
	sqlInsert := "insert into accounts (customer_id, opening_date, account_type, amount, status) values (? , ?, ?, ?, ?)"

	result, err := a.db.Exec(sqlInsert, account.CustomerId, account.OpeningDate, account.AccountType, account.Amount, account.Status)

	if err != nil {
		logger.Error("Unable to create account in database: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Unable to fetch last insert Id: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	account.AccountId = strconv.FormatInt(id, 10)
	return &account, nil
}

func (a AccountRepositoryMySql) insertTransaction(tx Transaction, dbTx *sql.Tx) (*Transaction, *errs.AppError) {
	sqlInsert := "insert into transactions (account_id, amount, transaction_type, transaction_date) values (?, ?, ?, ?)"
	result, err := a.db.Exec(sqlInsert, tx.AccountId, tx.Amount, tx.TransactionType, tx.TransactionDate)

	if err != nil {
		dbTx.Rollback()
		logger.Error("Unable to insert transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	txId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting transaction id: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	tx.TransactionId = strconv.FormatInt(txId, 10)
	return &tx, nil
}

func (a AccountRepositoryMySql) updateAmount(tx Transaction, dbTx *sql.Tx) *errs.AppError {
	var err error
	if tx.IsWithdrawal() {
		_, err = a.db.Exec("update accounts set amount = amount - ? where account_id = ?", tx.Amount, tx.AccountId)
	} else {
		_, err = a.db.Exec("update accounts set amount = amount + ? where account_id = ?", tx.Amount, tx.AccountId)
	}

	if err != nil {
		dbTx.Rollback()
		logger.Error("Error while updating account: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}
	return nil
}

func (a AccountRepositoryMySql) DoTransaction(tx Transaction) (*Transaction, *errs.AppError) {

	dbTx, err := a.db.Begin()
	if err != nil {
		logger.Error("Unable to begin DB transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	newTx, appError := a.insertTransaction(tx, dbTx)
	if appError != nil {
		return nil, appError
	}
	if appError = a.updateAmount(tx, dbTx); appError != nil {
		return nil, appError
	}

	err = dbTx.Commit()
	if err != nil {
		dbTx.Rollback()
		logger.Error("Error while committing transaction: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	account, appError := a.GetAccount(tx.AccountId)
	if appError != nil {
		return nil, appError
	}

	newTx.Amount = account.Amount
	return newTx, nil
}

func NewAccountRepositoryMySql(db *sqlx.DB) AccountRepositoryMySql {
	return AccountRepositoryMySql{db}
}
