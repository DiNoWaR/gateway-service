package server

import (
	"bytes"
	"encoding/json"
	model "github.com/dinowar/gateway-service/internal/pkg/domain"
	"log"
	"net/http"
)

var transactions = make(map[string]model.Transaction)

type Server struct {
}

func NewAppServer() Server {
	return Server{}
}

func (server *Server) Deposit(w http.ResponseWriter, r *http.Request) {
	var depositReq model.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&depositReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	transactionID := "tx_" + depositReq.AccountID // Генерация уникального ID транзакции

	transactions[transactionID] = model.Transaction{
		ID:       transactionID,
		UserID:   depositReq.AccountID,
		Status:   "pending", // Статус пока что "в ожидании"
		Amount:   depositReq.Amount,
		Currency: depositReq.Currency,
	}

	gatewayURL := determineGatewayURL(depositReq.Gateway)
	requestData := map[string]interface{}{
		"transaction_id": transactionID,
		"amount":         depositReq.Amount,
		"currency":       depositReq.Currency,
	}

	response, err := callPaymentGateway(gatewayURL, requestData)
	if err != nil {
		http.Error(w, "Error communicating with payment gateway", http.StatusInternalServerError)
		return
	}

	log.Printf("Payment gateway response for deposit: %v", response)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"transaction_id": transactionID})
}

func (server *Server) Withdraw(w http.ResponseWriter, r *http.Request) {
	var withdrawReq model.WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&withdrawReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	transactionID := "tx_" + withdrawReq.AccountID // Генерация уникального ID транзакции

	transactions[transactionID] = model.Transaction{
		ID:       transactionID,
		UserID:   withdrawReq.AccountID,
		Status:   "pending", // Статус пока что "в ожидании"
		Amount:   withdrawReq.Amount,
		Currency: withdrawReq.Currency,
	}

	gatewayURL := determineGatewayURL(withdrawReq.Gateway)
	requestData := map[string]interface{}{
		"transaction_id": transactionID,
		"amount":         withdrawReq.Amount,
		"currency":       withdrawReq.Currency,
	}

	response, err := callPaymentGateway(gatewayURL, requestData)
	if err != nil {
		http.Error(w, "Error communicating with payment gateway", http.StatusInternalServerError)
		return
	}

	log.Printf("Payment gateway response for withdrawal: %v", response)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"transaction_id": transactionID})
}

func (server *Server) Callback(w http.ResponseWriter, r *http.Request) {
	var callbackResponse model.CallbackResponse
	if err := json.NewDecoder(r.Body).Decode(&callbackResponse); err != nil {
		http.Error(w, "Invalid callback", http.StatusBadRequest)
		return
	}

	go processCallback(callbackResponse)

	w.WriteHeader(http.StatusOK)
}

func (server *Server) CheckTransactionStatus(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("transaction_id")
	transaction, exists := transactions[transactionID]

	if !exists {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(transaction)
}

func processCallback(callbackResponse model.CallbackResponse) {
	log.Printf("Received callback for transaction %s with status: %s", callbackResponse.TransactionID, callbackResponse.Status)

	transaction, exists := transactions[callbackResponse.TransactionID]
	if !exists {
		log.Printf("Transaction %s not found", callbackResponse.TransactionID)
		return
	}

	transaction.Status = callbackResponse.Status
	transactions[callbackResponse.TransactionID] = transaction

	sendEmailNotification(transaction.UserID, callbackResponse.Status)
}

func sendEmailNotification(userID string, status string) {
	log.Printf("Sending email to user %s: Your transaction status is: %s", userID, status)
}

func callPaymentGateway(gatewayURL string, requestData map[string]interface{}) (interface{}, error) {
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(gatewayURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func determineGatewayURL(gateway string) string {
	switch gateway {
	case "A":
		return "https://api.gateway-a.com/deposit"
	case "B":
		return "https://api.gateway-b.com/deposit"
	default:
		return ""
	}
}
