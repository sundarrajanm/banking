package domain

import (
	"banking/dto"
	"banking/errs"
)

type Customer struct {
	Id          string `db:"customer_id"`
	Name        string
	City        string
	Zipcode     string
	DateOfBirth string `db:"date_of_birth"`
	Status      string
}

func (c Customer) statusAsText() string {
	statusText := "inactive"
	if c.Status == "1" {
		statusText = "active"
	}
	return statusText
}

func (c Customer) ToDTO() dto.CustomerResponse {
	return dto.CustomerResponse{
		Id:          c.Id,
		Name:        c.Name,
		City:        c.City,
		DateOfBirth: c.DateOfBirth,
		Zipcode:     c.Zipcode,
		Status:      c.statusAsText(),
	}
}

type CustomerRepository interface {
	FindAllByStatus(string) ([]Customer, *errs.AppError)
	FindAll() ([]Customer, *errs.AppError)
	ById(string) (*Customer, *errs.AppError)
}
