package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dinowar/gateway-service/internal/pkg/config"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	"github.com/google/uuid"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

var (
	logger *zap.Logger
)

const (
	gatewayId = "rest_gateway"
)

func init() {
	logger, _ = zap.NewDevelopment()
}

func asyncProcessDeposit(referenceId string, callbackURL string) {
	// imitation of delay
	time.Sleep(1 * time.Second)

	callbackData := map[string]string{
		"transaction_id": uuid.NewString(),
		"reference_id":   referenceId,
		"status":         "SUCCESS",
		"message":        "deposit processed successfully",
	}

	sendCallback(callbackURL, callbackData)
}

func asyncProcessWithdraw(referenceId string, callbackURL string) {
	// imitation of delay
	time.Sleep(1 * time.Second)

	callbackData := map[string]string{
		"transaction_id": uuid.NewString(),
		"reference_id":   referenceId,
		"status":         "SUCCESS",
		"message":        "withdraw processed successfully",
	}

	sendCallback(callbackURL, callbackData)
}

func sendCallback(callbackURL string, data map[string]string) {
	reqBody, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		logger.Error("error marshalling callback data", zap.Error(marshalErr))
		return
	}

	resp, requestErr := http.Post(callbackURL, "application/json", bytes.NewBuffer(reqBody))
	if requestErr != nil {
		logger.Error("error calling callback url:", zap.String("url", callbackURL), zap.Error(requestErr))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Info("callback returned non-OK status", zap.String("status", resp.Status))
	}
}

func depositHandler(w http.ResponseWriter, r *http.Request) {
	var req DepositReq
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		logger.Error("error decoding request", zap.Error(decodeErr))
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	callbackURL := r.Header.Get("Callback-URL")
	if callbackURL == "" {
		http.Error(w, "missing Callback-URL header", http.StatusBadRequest)
		return
	}

	resp := DepositResponse{
		Gateway:       gatewayId,
		TransactionID: uuid.NewString(),
		AccountID:     req.AccountID,
		Status:        StatusPending,
		Message:       "deposit request is being processed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	go asyncProcessDeposit(req.ReferenceID, callbackURL)
}

func withdrawHandler(w http.ResponseWriter, r *http.Request) {
	var req WithdrawReq
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		logger.Error("error decoding request", zap.Error(decodeErr))
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	callbackURL := r.Header.Get("Callback-URL")
	if callbackURL == "" {
		http.Error(w, "missing Callback-URL header", http.StatusBadRequest)
		return
	}

	resp := WithdrawResponse{
		Gateway:       gatewayId,
		TransactionID: uuid.NewString(),
		AccountID:     req.AccountID,
		Status:        StatusPending,
		Message:       "withdraw request is being processed",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	go asyncProcessWithdraw(req.ReferenceID, callbackURL)
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
