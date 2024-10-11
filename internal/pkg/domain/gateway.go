package domain

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type PaymentGateway interface {
	ProcessDeposit(r *http.Request) (Transaction, error)
	ProcessWithdraw(r *http.Request) (Transaction, error)
}

type JsonGateway struct{}

func (gateway *JsonGateway) ProcessDeposit(r *http.Request) (Transaction, error) {
	var gatewayReq JsonGatewayRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&gatewayReq)
	if decodeErr != nil {
		return Transaction{}, fmt.Errorf("invalid request body: %v", decodeErr)
	}
	log.Printf("processing jsonGateway deposit: %+v", gatewayReq)
	resp := Transaction{
		ID:        uuid.New().String(),
		Amount:    gatewayReq.Amount,
		Currency:  gatewayReq.Currency,
		Status:    "pending",
		Operation: "deposit",
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	return resp, nil
}

func (gateway *JsonGateway) ProcessWithdraw(r *http.Request) (Transaction, error) {
	var gatewayReq JsonGatewayRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&gatewayReq)
	if decodeErr != nil {
		return Transaction{}, fmt.Errorf("invalid request body: %v", decodeErr)
	}
	log.Printf("processing jsonGateway deposit: %+v", gatewayReq)
	resp := Transaction{
		ID:        uuid.New().String(),
		Amount:    gatewayReq.Amount,
		Currency:  gatewayReq.Currency,
		Status:    "pending",
		Operation: "withdraw",
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	return resp, nil
}

type XMLGateway struct{}

func (gateway *XMLGateway) ProcessDeposit(r *http.Request) (Transaction, error) {
	return Transaction{}, fmt.Errorf("deposit not supported for Gateway B")
}

func (gateway *XMLGateway) ProcessWithdraw(r *http.Request) (Transaction, error) {
	var gatewayReq XMLGatewayRequest
	err := xml.NewDecoder(r.Body).Decode(&gatewayReq)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid request body: %v", err)
	}
	log.Printf("Processing Gateway B withdraw: %+v", gatewayReq)
	resp := Transaction{
		ID:        uuid.NewString(),
		Amount:    gatewayReq.Amount,
		Currency:  gatewayReq.Currency,
		Status:    "pending",
		Operation: "withdraw",
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	return resp, nil
}
