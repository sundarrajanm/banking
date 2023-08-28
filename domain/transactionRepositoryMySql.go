package domain

import (
	"banking/errs"
	"banking/logger"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type TransactionRepositoryMySql struct {
	db *sqlx.DB
}

func (t TransactionRepositoryMySql) Save(tx Transaction) (*Transaction, *errs.AppError) {
	sqlInsert := "insert into transactions (account_id, amount, transaction_type, transaction_date) values (?, ?, ?, ?)"

	result, err := t.db.Exec(sqlInsert, tx.AccountId, tx.Amount, tx.TransactionType, tx.TransactionDate)

	if err != nil {
		logger.Error("Unable to create transaction in database: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Unable to fetch last insert transaction Id: " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	tx.TransactionId = strconv.FormatInt(id, 10)
	return &tx, nil
}

func NewTransactionRepositoryMySql(db *sqlx.DB) TransactionRepositoryMySql {
	return TransactionRepositoryMySql{db}
}
