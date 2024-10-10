package main

import (
	"github.com/dinowar/gateway-service/internal/pkg/constants"
	"github.com/dinowar/gateway-service/internal/pkg/server"
	"log"
	"net/http"
)

func main() {
	appServer := server.NewAppServer()

	http.HandleFunc(constants.Deposit, appServer.Deposit)
	http.HandleFunc(constants.Withdraw, appServer.Withdraw)
	http.HandleFunc(constants.Callback, appServer.Callback)
	http.HandleFunc(constants.TransactionStatus, appServer.CheckTransactionStatus)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
