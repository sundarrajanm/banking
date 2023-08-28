package app

import (
	"banking/dto"
	"banking/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type TransactionHandler struct {
	service service.TransactionService
}

func (th *TransactionHandler) executeTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.NewTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	vars := mux.Vars(r)
	req.AccountId = vars["account_id"]
	txResponse, appError := th.service.NewTransaction(req)
	if appError != nil {
		writeResponse(w, appError.Code, appError.AsMessage())
		return
	}

	writeResponse(w, http.StatusOK, txResponse)
}
