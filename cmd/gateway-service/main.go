package main

import (
	"github.com/dinowar/gateway-service/internal/pkg/domain"
	"github.com/dinowar/gateway-service/internal/pkg/server"
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"github.com/dinowar/gateway-service/internal/pkg/util"
	"log"
	"net/http"
)

func main() {
	db, dbErr := util.InitDB()
	if dbErr != nil {
		log.Fatal("failed to initialize database", dbErr)
	}

	repService := service.NewRepositoryService(db)
	logService := service.NewLogService()
	appServer := server.NewAppServer(repService, logService)

	appServer.RegisterGateway("JsonGateway", &domain.JsonGateway{})
	appServer.RegisterGateway("XMLGateway", &domain.XMLGateway{})

	http.HandleFunc("/deposit", appServer.HandleDeposit)
	http.HandleFunc("/withdraw", appServer.HandleWithdraw)
	http.HandleFunc("/callback", appServer.HandleCallback)
	http.HandleFunc("/transaction", appServer.HandleGetTransaction)
	http.HandleFunc("/transactions", appServer.HandleGetAllTransactions)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
