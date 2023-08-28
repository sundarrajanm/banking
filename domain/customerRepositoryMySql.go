package domain

import (
	"database/sql"
	"log"
	"time"

	"banking/errs"

	_ "github.com/go-sql-driver/mysql"
)

type CustomerRepositoryMySql struct {
	db *sql.DB
}

func (s CustomerRepositoryMySql) FindAllByStatus(status string) ([]Customer, *errs.AppError) {
	findAllByStatusSql := "select * from customers where status = ?"
	rows, err := s.db.Query(findAllByStatusSql, status)
	return toCustomers(rows, err)
}

func toCustomers(rows *sql.Rows, err error) ([]Customer, *errs.AppError) {
	if err != nil {
		log.Println("Error while querying customer table: " + err.Error())
		return nil, errs.NewUnexpectedError("unexpected database error")
	}

	customers := make([]Customer, 0)
	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.DateOfBirth, &c.Status)

		if err != nil {
			log.Println("Error while scanning customers: " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (s CustomerRepositoryMySql) FindAll() ([]Customer, *errs.AppError) {
	findAllSql := "select * from customers"
	rows, err := s.db.Query(findAllSql)
	return toCustomers(rows, err)
}

func (s CustomerRepositoryMySql) ById(id string) (*Customer, *errs.AppError) {
	findById := "select * from customers where customer_id = ?"
	row := s.db.QueryRow(findById, id)
	var c Customer
	err := row.Scan(&c.Id, &c.Name, &c.City, &c.Zipcode, &c.DateOfBirth, &c.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("customer not found")
		} else {
			log.Println("Error while scanning customers: " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}
	return &c, nil
}

func NewCustomerRepositoryMySql() CustomerRepositoryMySql {
	db, err := sql.Open("mysql", "root:codecamp@/banking")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return CustomerRepositoryMySql{db}
}
