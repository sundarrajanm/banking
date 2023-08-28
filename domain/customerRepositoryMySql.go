package domain

import (
	"database/sql"
	"time"

	"banking/errs"
	"banking/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type CustomerRepositoryMySql struct {
	db *sqlx.DB
}

func (s CustomerRepositoryMySql) FindAllByStatus(status string) ([]Customer, *errs.AppError) {
	findAllByStatusSql := "select * from customers where status = ?"
	customers := make([]Customer, 0)
	err := s.db.Select(&customers, findAllByStatusSql, status)
	if err != nil {
		logger.Error("Error while scanning customers: " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}
	return customers, nil
}

func (s CustomerRepositoryMySql) FindAll() ([]Customer, *errs.AppError) {
	findAllSql := "select * from customers"
	customers := make([]Customer, 0)
	err := s.db.Select(&customers, findAllSql)
	if err != nil {
		logger.Error("Error while scanning customers: " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}
	return customers, nil
}

func (s CustomerRepositoryMySql) ById(id string) (*Customer, *errs.AppError) {
	findById := "select * from customers where customer_id = ?"
	var c Customer
	err := s.db.Get(&c, findById, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("customer not found")
		} else {
			logger.Error("Error while scanning customers: " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}
	return &c, nil
}

func NewCustomerRepositoryMySql() CustomerRepositoryMySql {
	db, err := sqlx.Open("mysql", "root:codecamp@/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return CustomerRepositoryMySql{db}
}
