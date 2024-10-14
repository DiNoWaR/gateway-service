package main

import (
	"context"
	"fmt"
	"github.com/dinowar/gateway-service/internal/pkg/config"
	"github.com/dinowar/gateway-service/internal/pkg/gateway"
	"github.com/dinowar/gateway-service/internal/pkg/server"
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"github.com/dinowar/gateway-service/internal/pkg/util"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

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
	logService := service.NewLogService(logger)
	appServer := server.NewAppServer(repService, logService, serviceConfig)

	// registering gateways
	appServer.RegisterGateway(serviceConfig.RestGatewayConfig.GatewayId,
		&gateway.RestGateway{
			BaseURL: fmt.Sprintf("%s:%s", serviceConfig.RestGatewayConfig.Host, serviceConfig.RestGatewayConfig.Port),
			Logger:  logger,
		})

	appServer.RegisterGateway(serviceConfig.SoapGatewayConfig.GatewayId,
		&gateway.SoapGateway{
			Endpoint: fmt.Sprintf("%s:%s%s", serviceConfig.SoapGatewayConfig.EndpointHost, serviceConfig.SoapGatewayConfig.EndpointPort, serviceConfig.SoapGatewayConfig.Endpoint),
			Logger:   logger,
		})

	http.HandleFunc("/deposit", appServer.HandleDeposit)
	http.HandleFunc("/withdraw", appServer.HandleWithdraw)
	http.HandleFunc("/callback", appServer.HandleCallback)
	http.HandleFunc("/transaction", appServer.HandleGetTransaction)
	http.HandleFunc("/transactions", appServer.HandleGetTransactions)

	log.Println(fmt.Sprintf("service started on port: %s", serviceConfig.ServicePort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serviceConfig.ServicePort), nil))
}
