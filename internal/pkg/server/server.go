package server

import (
	"context"
	"encoding/json"
	"fmt"
	. "github.com/dinowar/gateway-service/internal/pkg/domain"
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	gateways = map[string]PaymentGateway{}
)

type Server struct {
	rep        *service.RepositoryService
	logger     *service.LogService
	mux        sync.RWMutex
	workerPool chan struct{}
}

func NewAppServer(rep *service.RepositoryService, logger *service.LogService, workerCount int) *Server {
	return &Server{
		rep:        rep,
		logger:     logger,
		workerPool: make(chan struct{}, workerCount),
	}
}

func (server *Server) RegisterGateway(name string, gateway PaymentGateway) {
	server.mux.Lock()
	defer server.mux.Unlock()
	gateways[name] = gateway
	log.Printf("registered gateway: %s", name)
}

func (server *Server) handleTransaction(w http.ResponseWriter, r *http.Request, transactionType string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		AccountID string  `json:"accountId"`
		Amount    float64 `json:"amount"`
		Currency  string  `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AccountID == "" || req.Amount <= 0 {
		http.Error(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	gateway := r.URL.Query().Get("gateway")
	if gateway == "" {
		http.Error(w, "Missing gateway parameter", http.StatusBadRequest)
		return
	}

	if gateway == "" {
		http.Error(w, "Missing gateway parameter", http.StatusBadRequest)
		return
	}

	server.mux.RLock()
	gw, exists := gateways[gateway]
	server.mux.RUnlock()
	if !exists {
		http.Error(w, "Gateway not found", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	select {
	case server.workerPool <- struct{}{}:
		defer func() { <-server.workerPool }()
		var resp Transaction
		var err error
		if transactionType == "deposit" {
			resp, err = gw.ProcessDeposit(req)
		} else {
			resp, err = gw.ProcessWithdraw(req)
		}
		if err != nil {
			log.Printf("Failed to process %s: %v", transactionType, err)
			http.Error(w, "Failed to process transaction", http.StatusInternalServerError)
			return
		}

		transactionID := fmt.Sprintf("txn-%d", time.Now().UnixNano())
		resp.ID = transactionID
		if err := server.rep.SaveTransaction(resp); err != nil {
			http.Error(w, "Failed to save transaction", http.StatusInternalServerError)
			return
		}

		server.logger.LogTransaction(resp)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
	}
}

func (server *Server) HandleDeposit(w http.ResponseWriter, r *http.Request) {
	server.handleTransaction(w, r, "$2")
}

func (server *Server) HandleWithdraw(w http.ResponseWriter, r *http.Request) {
	server.handleTransaction(w, r, "$2")
}

func (server *Server) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received callback: %s", string(body))
	transactionID := r.URL.Query().Get("transaction_id")
	if transactionID == "" {
		http.Error(w, "Missing transaction_id", http.StatusBadRequest)
		return
	}

	txn, trxGetErr := server.rep.GetTransaction(transactionID)
	if trxGetErr != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}
	txn.Status = "success"

	if trxSaveErr := server.rep.SaveTransaction(txn); trxSaveErr != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	server.logger.LogTransaction(txn)
	w.WriteHeader(http.StatusOK)
}

func (server *Server) HandleGetTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	transactionID := r.URL.Query().Get("transaction_id")
	if transactionID == "" {
		http.Error(w, "Missing transaction_id", http.StatusBadRequest)
		return
	}

	txn, err := server.rep.GetTransaction(transactionID)
	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txn)
}

func (server *Server) HandleGetAllTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	accountId := r.URL.Query().Get("accountId")

	page := r.URL.Query().Get("page")
	pageSize := r.URL.Query().Get("page_size")
	transactions, transactionsErr := server.rep.GetTransactions(accountId, page, pageSize)
	if transactionsErr != nil {
		http.Error(w, "Transactions not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
