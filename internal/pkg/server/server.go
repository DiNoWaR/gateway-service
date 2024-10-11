package server

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	. "github.com/dinowar/gateway-service/internal/pkg/domain"
	"github.com/dinowar/gateway-service/internal/pkg/service"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	transactionsMutex = &sync.Mutex{}
	gateways          = map[string]PaymentGateway{}
)

type Server struct {
	rep    *service.RepositoryService
	logger *service.LogService
}

func NewAppServer(rep *service.RepositoryService, logger *service.LogService) *Server {
	return &Server{rep: rep, logger: logger}
}

func (server *Server) RegisterGateway(name string, gateway PaymentGateway) {
	gateways[name] = gateway
	log.Printf("registered gateway: %s", name)
}

func (server *Server) HandleDeposit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	gateway := r.URL.Query().Get("gateway")
	if gateway == "" {
		http.Error(w, "Missing gateway parameter", http.StatusBadRequest)
		return
	}

	gw, exists := gateways[gateway]
	if !exists {
		http.Error(w, "Gateway not found", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ch := make(chan Transaction, 1)

	go func() {
		resp, err := gw.ProcessDeposit(r)
		if err != nil {
			log.Printf("Failed to process deposit: %v", err)
			return
		}
		ch <- resp
	}()

	select {
	case transaction := <-ch:
		transactionID := fmt.Sprintf("txn-%d", time.Now().UnixNano())
		transaction.ID = transactionID
		transactionsMutex.Lock()
		err := server.rep.SaveTransaction(transaction)
		transactionsMutex.Unlock()
		if err != nil {
			http.Error(w, "Failed to save transaction", http.StatusInternalServerError)
			return
		}

		server.logger.LogTransaction(transaction)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transaction)
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
	}
}

func (server *Server) HandleWithdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	gateway := r.URL.Query().Get("gateway")
	if gateway == "" {
		http.Error(w, "Missing gateway parameter", http.StatusBadRequest)
		return
	}

	gw, exists := gateways[gateway]
	if !exists {
		http.Error(w, "Gateway not found", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ch := make(chan Transaction, 1)

	go func() {
		resp, err := gw.ProcessWithdraw(r)
		if err != nil {
			log.Printf("Failed to process withdraw: %v", err)
			return
		}
		ch <- resp
	}()

	select {
	case transaction := <-ch:
		transactionID := fmt.Sprintf("txn-%d", time.Now().UnixNano())
		transaction.ID = transactionID
		transactionsMutex.Lock()
		trxErr := server.rep.SaveTransaction(transaction)
		transactionsMutex.Unlock()
		if trxErr != nil {
			http.Error(w, "Failed to save transaction", http.StatusInternalServerError)
			return
		}

		server.logger.LogTransaction(transaction)
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(transaction)
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusRequestTimeout)
	}
}

func (server *Server) HandleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
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

	transactionsMutex.Lock()
	txn, trxGetErr := server.rep.GetTransaction(transactionID)
	if trxGetErr != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		transactionsMutex.Unlock()
		return
	}
	txn.Status = "updated"

	trxSaveErr := server.rep.SaveTransaction(txn)
	transactionsMutex.Unlock()
	if trxSaveErr != nil {
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

	transactionsMutex.Lock()
	txn, err := server.rep.GetTransaction(transactionID)
	transactionsMutex.Unlock()
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
	if accountId == "" {
		http.Error(w, "Missing transaction_id", http.StatusBadRequest)
		return
	}

	transactions, transactionsErr := server.rep.GetTransactions(accountId)
	if transactionsErr != nil {
		http.Error(w, "Transactions not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
