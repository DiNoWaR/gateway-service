package server

import (
	"encoding/json"
	"fmt"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	gateways "github.com/dinowar/gateway-service/internal/pkg/gateway"
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"net/http"
	"time"
)

type Server struct {
	rep      *service.RepositoryService
	logger   *service.LogService
	gateways map[string]gateways.PaymentGateway
}

func NewAppServer(rep *service.RepositoryService, logger *service.LogService) *Server {
	return &Server{
		rep:      rep,
		logger:   logger,
		gateways: make(map[string]gateways.PaymentGateway),
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

	var req ClientDepositRequest
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
		Amount:      decimal.NewFromFloat(req.Amount),
		Currency:    req.Currency,
		Status:      StatusPending,
		Operation:   Deposit,
		Ts:          time.Now(),
	}

	//trxErr := server.rep.SaveTransaction(txn)
	//if trxErr != nil {
	//	http.Error(w, "error saving transaction", http.StatusInternalServerError)
	//	server.logger.LogError("HandleDeposit: error saving transaction: %v", trxErr)
	//	return
	//}

	depositReq := gateways.DepositReq{
		Amount:      req.Amount,
		Currency:    req.Currency,
		ReferenceID: referenceId,
		AccountID:   req.AccountID,
	}

	depositResp, gatewayErr := gateway.ProcessDeposit(depositReq)
	if gatewayErr != nil {
		http.Error(w, "Error processing deposit", http.StatusInternalServerError)
		server.logger.LogError("HandleDeposit: error processing deposit: %v", gatewayErr)
		txn.Status = StatusFailed
		//_ = server.rep.UpdateTransaction(txn)
		return
	}

	txn.Status = depositResp.Status
	txn.Id = depositResp.TransactionID
	//err = server.rep.UpdateTransaction(txn)
	//if err != nil {
	//	http.Error(w, "Error updating transaction", http.StatusInternalServerError)
	//	server.logger.LogError("HandleDeposit: error updating transaction: %v", err)
	//	return
	//}

	response := map[string]interface{}{
		"transaction_id": txn.Id,
		"status":         txn.Status,
		"message":        depositResp.Message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoderErr := json.NewEncoder(w).Encode(response)
	if encoderErr != nil {
		server.logger.LogError("HandleDeposit: error encoding response: %v", encoderErr)
	}
}

func (server *Server) HandleWithdraw(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) HandleCallback(w http.ResponseWriter, r *http.Request) {

}

func (server *Server) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {

}
