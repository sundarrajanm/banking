package app

import (
	"banking/dto"
	"banking/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type AccountHandler struct {
	service service.AccountService
}

func (ah AccountHandler) NewAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.NewAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	vars := mux.Vars(r)
	req.CustomerId = vars["customer_id"]
	account, appError := ah.service.NewAccount(req)
	if appError != nil {
		writeResponse(w, appError.Code, appError.AsMessage())
		return
	}

	writeResponse(w, http.StatusCreated, account)
}

func (ah AccountHandler) executeTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.NewTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	vars := mux.Vars(r)
	req.AccountId = vars["account_id"]
	txResponse, appError := ah.service.NewTransaction(req)
	if appError != nil {
		writeResponse(w, appError.Code, appError.AsMessage())
		return
	}

	writeResponse(w, http.StatusOK, txResponse)
}
