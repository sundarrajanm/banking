package app

import (
	"banking/dto"
	"banking/errs"
	"banking/mocks/service"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var router *mux.Router
var ch CustomerHandler
var mockService *service.MockCustomerService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockCustomerService(ctrl)
	ch = CustomerHandler{mockService}
	router = mux.NewRouter()
	router.HandleFunc("/customers", ch.getAllCustomers)
	return func() {
		router = nil
		defer ctrl.Finish()
	}
}

func Test_should_return_customers_with_status_code_200(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	dummyCustomers := []dto.CustomerResponse{
		{"1001", "Ashish", "New Delhi", "110075", "2000-01-01", "1"},
		{"1002", "Rob", "New Delhi", "110075", "2000-01-01", "1"},
	}
	mockService.EXPECT().GetAllCustomers().Return(dummyCustomers, nil)
	req, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing status code")
	}
}

func Test_should_return_status_code_500_with_error_message(t *testing.T) {
	tearDown := setup(t)
	defer tearDown()

	mockService.EXPECT().GetAllCustomers().Return(nil, errs.NewUnexpectedError("some database error"))
	req, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing status code")
	}
}
