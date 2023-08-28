package service

import (
	"banking/domain"
	"banking/errs"
)

type CustomerService interface {
	GetAllCustomersByStatus(string) ([]domain.Customer, *errs.AppError)
	GetAllCustomers() ([]domain.Customer, *errs.AppError)
	GetCustomer(string) (*domain.Customer, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomersByStatus(status string) ([]domain.Customer, *errs.AppError) {
	statusValue := "-1"
	if status == "active" {
		statusValue = "1"
	} else if status == "inactive" {
		statusValue = "0"
	}

	return s.repo.FindAllByStatus(statusValue)
}

func (s DefaultCustomerService) GetAllCustomers() ([]domain.Customer, *errs.AppError) {
	return s.repo.FindAll()
}

func (s DefaultCustomerService) GetCustomer(id string) (*domain.Customer, *errs.AppError) {
	return s.repo.ById(id)
}

func NewCustomerService(repo domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repo}
}
