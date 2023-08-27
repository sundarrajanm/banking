package app

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"banking/service"

	"github.com/gorilla/mux"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

type Customer struct {
	Name    string `json:"full_name" xml:"name"`
	City    string `json:"city" xml:"city"`
	Zipcode string `json:"zip_code" xml:"zipcode"`
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

type CustomerHandler struct {
	service service.CustomerService
}

func (ch *CustomerHandler) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers, _ := ch.service.GetAllCustomers()

	if r.Header.Get("Content-Type") == "application/xml" {
		w.Header().Add("content-type", "application/xml")
		xml.NewEncoder(w).Encode(customers)
	} else {
		w.Header().Add("content-type", "application/json")
		json.NewEncoder(w).Encode(customers)
	}
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprint(w, vars["customer_id"])
}
