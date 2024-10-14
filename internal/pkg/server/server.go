package server

import (
	"encoding/json"
	"fmt"
	"github.com/dinowar/gateway-service/internal/pkg/config"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	gateways "github.com/dinowar/gateway-service/internal/pkg/gateway"
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"time"
)

type Server struct {
	rep      *service.RepositoryService
	logger   *service.LogService
	gateways map[string]gateways.PaymentGateway
	config   *config.ServiceConfig
}

func NewAppServer(rep *service.RepositoryService, logger *service.LogService, config *config.ServiceConfig) *Server {
	return &Server{
		rep:      rep,
		logger:   logger,
		gateways: make(map[string]gateways.PaymentGateway),
		config:   config,
	}
}

func (server *Server) RegisterGateway(gatewayId string, gateway gateways.PaymentGateway) {
	server.gateways[gatewayId] = gateway
}

func (server *Server) HandleDeposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ClientRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		server.logger.LogError("HandleDeposit: error decoding request body: %v", decodeErr)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 || req.Currency == "" || req.AccountID == "" || req.GatewayID == "" {
		http.Error(w, "missing or invalid fields in request body", http.StatusBadRequest)
		return
	}

	gateway, exists := server.gateways[req.GatewayID]
	if !exists {
		http.Error(w, fmt.Sprintf("gateway %s not found", req.GatewayID), http.StatusBadRequest)
		return
	}

	referenceId := uuid.NewString()
	txn := &Transaction{
		ReferenceId: referenceId,
		AccountId:   req.AccountID,
		GatewayId:   req.GatewayID,
		Amount:      decimal.NewFromFloat(req.Amount),
		Currency:    req.Currency,
		Status:      StatusPending,
		Operation:   Deposit,
		Ts:          time.Now(),
	}

	depositReq := DepositReq{
		Amount:      req.Amount,
		Currency:    req.Currency,
		ReferenceID: referenceId,
		AccountID:   req.AccountID,
	}

	depositResp, gatewayErr := gateway.ProcessDeposit(depositReq, server.config.ServiceCallbackEndpoint)
	if gatewayErr != nil {
		http.Error(w, "error processing deposit", http.StatusInternalServerError)
		server.logger.LogError("HandleDeposit: error processing deposit: %v", gatewayErr)
		return
	}

	trxErr := server.rep.SaveTransaction(txn)
	if trxErr != nil {
		http.Error(w, "error saving transaction", http.StatusInternalServerError)
		server.logger.LogError("HandleDeposit: error saving transaction: %v", trxErr)
		return
	}

	response := map[string]interface{}{
		"transaction_status": txn.Status,
		"operation_type":     Deposit,
		"gateway":            depositResp.Gateway,
		"account_id":         req.AccountID,
		"reference_id":       txn.ReferenceId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoderErr := json.NewEncoder(w).Encode(response)
	if encoderErr != nil {
		server.logger.LogError("HandleDeposit: error encoding response: %v", encoderErr)
	}
}

func (server *Server) HandleWithdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ClientRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		server.logger.LogError("HandleDeposit: error decoding request body: %v", decodeErr)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 || req.Currency == "" || req.AccountID == "" || req.GatewayID == "" {
		http.Error(w, "missing or invalid fields in request body", http.StatusBadRequest)
		return
	}

	gateway, exists := server.gateways[req.GatewayID]
	if !exists {
		http.Error(w, fmt.Sprintf("gateway %s not found", req.GatewayID), http.StatusBadRequest)
		return
	}

	referenceId := uuid.NewString()
	txn := &Transaction{
		ReferenceId: referenceId,
		AccountId:   req.AccountID,
		GatewayId:   req.GatewayID,
		Amount:      decimal.NewFromFloat(req.Amount),
		Currency:    req.Currency,
		Status:      StatusPending,
		Operation:   Withdraw,
		Ts:          time.Now(),
	}

	withdrawReq := WithdrawReq{
		Amount:      req.Amount,
		Currency:    req.Currency,
		ReferenceID: referenceId,
		AccountID:   req.AccountID,
	}

	depositResp, gatewayErr := gateway.ProcessWithdrawal(withdrawReq, server.config.ServiceCallbackEndpoint)
	if gatewayErr != nil {
		http.Error(w, "error processing deposit", http.StatusInternalServerError)
		server.logger.LogError("HandleDeposit: error processing deposit: %v", gatewayErr)
		return
	}

	trxErr := server.rep.SaveTransaction(txn)
	if trxErr != nil {
		http.Error(w, "error saving transaction", http.StatusInternalServerError)
		server.logger.LogError("HandleDeposit: error saving transaction: %v", trxErr)
		return
	}

	response := map[string]interface{}{
		"transaction_status": txn.Status,
		"operation_type":     Deposit,
		"gateway":            depositResp.Gateway,
		"account_id":         req.AccountID,
		"reference_id":       txn.ReferenceId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoderErr := json.NewEncoder(w).Encode(response)
	if encoderErr != nil {
		server.logger.LogError("HandleDeposit: error encoding response: %v", encoderErr)
	}
}

func (server *Server) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		server.logger.LogError("HandleCallback: error reading request body: %v", readErr)
		http.Error(w, "error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req CallbackPayload
	decodeErr := json.Unmarshal(body, &req)
	if decodeErr != nil {
		server.logger.LogError("HandleCallback: error decoding request body: %v", decodeErr)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	trxErr := server.rep.UpdateTransaction(&Transaction{
		Id:          req.TransactionId,
		ReferenceId: req.ReferenceId,
		Status:      TransactionStatus(req.Status),
		Message:     req.Message,
	})
	if trxErr != nil {
		http.Error(w, "error updating transaction", http.StatusInternalServerError)
		server.logger.LogError("HandleCallback: error saving transaction: %v", trxErr)
		return
	}
	fmt.Println(string(body))
}

func (server *Server) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var transactionReq GetTransactionRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&transactionReq)
	if decodeErr != nil {
		server.logger.LogError("HandleGetTransaction: error decoding request body: %v", decodeErr)
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	if transactionReq.ReferenceId == "" {
		http.Error(w, "missing or invalid fields in request body", http.StatusBadRequest)
	}

	transaction, trErr := server.rep.GetTransaction(transactionReq.ReferenceId)
	if trErr != nil {
		server.logger.LogError("HandleGetTransaction: error getting transaction: %v", trErr)
		http.Error(w, "error getting transaction", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoderErr := json.NewEncoder(w).Encode(transaction)
	if encoderErr != nil {
		server.logger.LogError("HandleGetTransaction: error encoding response: %v", encoderErr)
	}
}

func (server *Server) HandleGetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var transactionReq GetTransactionsRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&transactionReq)
	if decodeErr != nil {
		server.logger.LogError("HandleGetTransaction: error decoding request body: %v", decodeErr)
		http.Error(w, "invalid request body", http.StatusBadRequest)
	}

	if transactionReq.AccountId == "" {
		http.Error(w, "missing or invalid fields in request body", http.StatusBadRequest)
	}

	transaction, trErr := server.rep.GetTransactions(transactionReq.AccountId)
	if trErr != nil {
		server.logger.LogError("HandleGetTransaction: error getting transaction: %v", trErr)
		http.Error(w, "error getting transaction", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoderErr := json.NewEncoder(w).Encode(transaction)
	if encoderErr != nil {
		server.logger.LogError("HandleGetTransaction: error encoding response: %v", encoderErr)
	}
}
