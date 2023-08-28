package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"banking/domain"
	"banking/errs"
	"banking/service"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	service service.CustomerService
}

func (ch *CustomerHandler) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	var err *errs.AppError
	var customers []domain.Customer

	status := r.URL.Query().Get("status")
	if status == "" {
		customers, err = ch.service.GetAllCustomers()
	} else {
		customers, err = ch.service.GetAllCustomersByStatus(status)
	}

	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
		return
	}

	writeResponse(w, http.StatusOK, customers)
}

func (ch *CustomerHandler) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	customer_id := vars["customer_id"]
	customer, err := ch.service.GetCustomer(customer_id)

	if err != nil {
		writeResponse(w, err.Code, err.AsMessage())
		return
	}

	writeResponse(w, http.StatusOK, customer)
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

///// /api/time assignment //////

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func GetTimeFromTimezone(tz string) (string, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return "", errors.New("invalid timezone")
	}
	return time.Now().In(loc).String(), nil
}

func handleTimezone(timezone string, w http.ResponseWriter) {
	response := map[string]string{}
	timezones := strings.Split(timezone, ",")

	for i := 0; i < len(timezones); i++ {
		tz := timezones[i]

		loc, err := GetTimeFromTimezone(strings.TrimSpace(tz))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, err.Error())
			return
		}

		response[tz] = loc
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleUTC(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"current_time": time.Now().UTC().String(),
	})
}
func getTime(w http.ResponseWriter, r *http.Request) {
	timezone := r.URL.Query().Get("tz")

	if timezone != "" {
		handleTimezone(timezone, w)
	} else {
		handleUTC(w)
	}
}

///// /api/time assignment //////
