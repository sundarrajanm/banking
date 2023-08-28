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

func (a AccountRepositoryMySql) Debit(amount float64, accountId string) *errs.AppError {
	WithdrawAmountSql := "update accounts set amount = amount - ? where account_id = ?"

	_, err := a.db.Exec(WithdrawAmountSql, amount, accountId)
	if err != nil {
		logger.Error("Unexpected database error while withdrawing amount: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}
	return nil
}

func (a AccountRepositoryMySql) Credit(amount float64, accountId string) *errs.AppError {
	DepositAmountSql := "update accounts set amount = amount + ? where account_id = ?"

	_, err := a.db.Exec(DepositAmountSql, amount, accountId)
	if err != nil {
		logger.Error("Unexpected database error while depositing amount: " + err.Error())
		return errs.NewUnexpectedError("Unexpected database error")
	}
	return nil
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

func NewAccountRepositoryMySql(db *sqlx.DB) AccountRepositoryMySql {
	return AccountRepositoryMySql{db}
}
