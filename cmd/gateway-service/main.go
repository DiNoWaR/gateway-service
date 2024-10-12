package main

import (
	"context"
	"github.com/dinowar/gateway-service/internal/pkg/config"
	"github.com/dinowar/gateway-service/internal/pkg/server"
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"github.com/dinowar/gateway-service/internal/pkg/util"
	"github.com/sethvargo/go-envconfig"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	serviceConfig := &config.ServiceConfig{}
	if configErr := envconfig.Process(ctx, serviceConfig); configErr != nil {
		log.Fatal(ctx, "failed to init config", configErr)
	}
	db, dbErr := util.InitDB(
		serviceConfig.DBConfig.Host,
		serviceConfig.DBConfig.Port,
		serviceConfig.DBConfig.Database,
		serviceConfig.DBConfig.Username,
		serviceConfig.DBConfig.Password,
	)
	if dbErr != nil {
		log.Println("failed to initialize database", dbErr)
	}

	repService := service.NewRepositoryService(db)
	logService := service.NewLogService()
	appServer := server.NewAppServer(repService, logService)

	http.HandleFunc("/deposit", appServer.HandleDeposit)
	http.HandleFunc("/withdraw", appServer.HandleWithdraw)
	http.HandleFunc("/callback", appServer.HandleCallback)
	http.HandleFunc("/transaction", appServer.HandleGetTransaction)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
