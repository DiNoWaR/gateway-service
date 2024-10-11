package domain

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
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
	err := json.NewDecoder(r.Body).Decode(&gatewayReq)
	if err != nil {
		return Transaction{}, fmt.Errorf("invalid request body: %v", err)
	}
	log.Printf("Processing Gateway A deposit: %+v", gatewayReq)
	resp := Transaction{
		ID:       "",
		Amount:   gatewayReq.Amount,
		Currency: gatewayReq.Currency,
		Status:   "success",
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	return resp, nil
}

func (gateway *JsonGateway) ProcessWithdraw(r *http.Request) (Transaction, error) {
	return Transaction{}, fmt.Errorf("withdraw not supported for Gateway A")
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
		ID:       "",
		Amount:   gatewayReq.Amount,
		Currency: gatewayReq.Currency,
		Status:   "success",
	}
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	return resp, nil
}
