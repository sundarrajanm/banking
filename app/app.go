package app

import (
	"banking/domain"
	"banking/logger"
	"banking/service"
	"errors"
	"os"
	"time"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
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

	db := getDBClient()
	acctRepo := domain.NewAccountRepositoryMySql(db)
	ch := CustomerHandler{service.NewCustomerService(domain.NewCustomerRepositoryMySql(db))}
	ah := AccountHandler{service.NewAccountService(acctRepo)}

	router := mux.NewRouter()

	router.
		HandleFunc("/customers", ch.getAllCustomers).
		Methods(http.MethodGet).
		Name("GetAllCustomers")

	router.
		HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer).
		Methods(http.MethodGet).
		Name("GetCustomer")

	router.
		HandleFunc("/customers/{customer_id:[0-9]+}/accounts", ah.NewAccount).
		Methods(http.MethodPost).
		Name("NewAccount")

	router.
		HandleFunc("/customers/{customer_id:[0-9]+}/accounts/{account_id:[0-9]+}/transaction", ah.executeTransaction).
		Methods(http.MethodPost).
		Name("NewTransaction")

	router.HandleFunc("/greet", greet).Methods(http.MethodGet)
	router.HandleFunc("/api/time", getTime)

	am := AuthMiddleware{domain.NewAuthRepository()}
	router.Use(am.authorizationHandler())

	serverAddr := os.Getenv("SERVER_ADDR")
	serverPort := os.Getenv("SERVER_PORT")
	serverUrl := fmt.Sprintf("%s:%s", serverAddr, serverPort)
	logger.Fatal(http.ListenAndServe(serverUrl, router))
}

func getDBClient() *sqlx.DB {
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPasswd := os.Getenv("DB_PASSWD")
	dbName := os.Getenv("DB_NAME")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPasswd, dbAddr, dbPort, dbName)

	db, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
