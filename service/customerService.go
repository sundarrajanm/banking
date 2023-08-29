package service

import (
	"banking/domain"
	"banking/dto"
	"banking/errs"
)

//go:generate mockgen -destination=../mocks/service/mockCustomerService.go -package=service banking/service CustomerService
type CustomerService interface {
	GetAllCustomersByStatus(string) ([]dto.CustomerResponse, *errs.AppError)
	GetAllCustomers() ([]dto.CustomerResponse, *errs.AppError)
	GetCustomer(string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (s DefaultCustomerService) GetAllCustomersByStatus(status string) ([]dto.CustomerResponse, *errs.AppError) {
	statusValue := "-1"
	if status == "active" {
		statusValue = "1"
	} else if status == "inactive" {
		statusValue = "0"
	}

	customers, appError := s.repo.FindAllByStatus(statusValue)
	if appError != nil {
		return nil, appError
	}

	return toDTOs(customers), nil
}

func toDTOs(customers []domain.Customer) []dto.CustomerResponse {
	custResponse := make([]dto.CustomerResponse, 0)
	for _, c := range customers {
		custResponse = append(custResponse, c.ToDTO())
	}
	return custResponse
}

func (s DefaultCustomerService) GetAllCustomers() ([]dto.CustomerResponse, *errs.AppError) {
	customers, appError := s.repo.FindAll()
	if appError != nil {
		return nil, appError
	}

	return toDTOs(customers), nil
}

func (s DefaultCustomerService) GetCustomer(id string) (*dto.CustomerResponse, *errs.AppError) {
	c, err := s.repo.ById(id)
	if err != nil {
		return nil, err
	}

	response := c.ToDTO()
	return &response, nil
}

func NewCustomerService(repo domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repo}
}
