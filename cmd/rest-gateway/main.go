package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dinowar/gateway-service/internal/pkg/config"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/common"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/rest_gateway"
	"github.com/google/uuid"
	"github.com/sethvargo/go-envconfig"
	"log"
	"net/http"
)

const (
	gatewayId = "rest_gateway"
)

func depositHandler(w http.ResponseWriter, r *http.Request) {
	var req DepositReq
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("received depositrequest: %+v\n", req)
	resp := DepositResponse{
		Gateway:       gatewayId,
		TransactionID: uuid.NewString(),
		AccountID:     req.AccountID,
		Status:        StatusSuccess,
		Message:       "deposit processed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func withdrawHandler(w http.ResponseWriter, r *http.Request) {
	var req WithdrawReq
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("received withdrawrequest: %+v\n", req)
	resp := WithdrawResponse{
		Gateway:       gatewayId,
		TransactionID: uuid.NewString(),
		AccountID:     req.AccountID,
		Status:        StatusSuccess,
		Message:       "withdrawal processed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	ctx := context.Background()
	serviceConfig := &config.ServiceConfig{}
	if configErr := envconfig.Process(ctx, serviceConfig); configErr != nil {
		log.Fatal(ctx, "failed to init config", configErr)
	}

	http.HandleFunc("/deposit", depositHandler)
	http.HandleFunc("/withdraw", withdrawHandler)

	address := fmt.Sprintf("%s:%s", serviceConfig.RestGatewayConfig.Host, serviceConfig.RestGatewayConfig.Port)

	fmt.Println(fmt.Sprintf("rest mock server running on address %s...", address))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serviceConfig.RestGatewayConfig.Port), nil))
}
