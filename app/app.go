package app

import (
	"banking/domain"
	"banking/logger"
	"banking/service"
	"errors"
	"os"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func EnvVarsNotEmptyCheck(vars ...string) {
	for _, envVar := range vars {
		if os.Getenv(envVar) == "" {
			logger.Fatal(errors.New("Mandatory configuration variable not found: " + envVar))
		}
	}
}

func Start() {

	EnvVarsNotEmptyCheck(
		"SERVER_ADDR", "SERVER_PORT",
		"DB_ADDR", "DB_PORT", "DB_USER", "DB_PASSWD", "DB_NAME",
	)

	router := mux.NewRouter()

	// ch := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryStub())}
	ch := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryMySql())}

	router.HandleFunc("/greet", greet).Methods(http.MethodGet)
	router.HandleFunc("/customers", ch.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers", createCustomer).Methods(http.MethodPost)
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer)
	router.HandleFunc("/api/time", getTime)

	serverAddr := os.Getenv("SERVER_ADDR")
	serverPort := os.Getenv("SERVER_PORT")
	serverUrl := fmt.Sprintf("%s:%s", serverAddr, serverPort)
	logger.Fatal(http.ListenAndServe(serverUrl, router))
}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Post request received")
}
