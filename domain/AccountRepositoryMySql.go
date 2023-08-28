package domain

import (
	"banking/errs"
	"banking/logger"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type AccountRepositoryMySql struct {
	db *sqlx.DB
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
