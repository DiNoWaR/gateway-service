package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/dinowar/gateway-service/internal/pkg/config"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	. "github.com/dinowar/gateway-service/internal/pkg/gateway"
	"github.com/google/uuid"
	"github.com/sethvargo/go-envconfig"
	"io"
	"log"
	"net/http"
)

const (
	gatewayId = "soap"
	soapNS    = "http://schemas.xmlsoap.org/soap/envelope/"
)

func soapHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		http.Error(w, readErr.Error(), http.StatusBadRequest)
		return
	}

	var envelope Envelope
	unmarshalErr := xml.Unmarshal(bodyBytes, &envelope)
	if unmarshalErr != nil {
		http.Error(w, unmarshalErr.Error(), http.StatusBadRequest)
		return
	}

	var response interface{}

	switch content := envelope.Body.Content.(type) {
	case DepositReq:
		req := content
		log.Printf("processing DepositRequest: %+v\n", req)
		response = DepositResponse{
			XMLName:       xml.Name{Local: "DepositResponse"},
			Gateway:       gatewayId,
			TransactionID: uuid.NewString(),
			Status:        StatusSuccess,
			Message:       "Deposit processed successfully",
			AccountID:     req.AccountID,
		}
	case WithdrawReq:
		req := content
		log.Printf("processing WithdrawRequest: %+v\n", req)
		response = WithdrawResponse{
			XMLName:       xml.Name{Local: "WithdrawResponse"},
			Gateway:       gatewayId,
			TransactionID: uuid.NewString(),
			Status:        StatusSuccess,
			Message:       "Withdrawal processed successfully",
			AccountID:     req.AccountID,
		}
	default:
		http.Error(w, "Unknown request", http.StatusBadRequest)
		return
	}

	sendSoapResponse(w, response)
}

func sendSoapResponse(w http.ResponseWriter, response interface{}) {
	soapResponse := Envelope{
		SoapNS: soapNS,
		Body:   Body{Content: response},
	}

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	encodeErr := xml.NewEncoder(w).Encode(soapResponse)
	if encodeErr != nil {
		log.Printf("error encoding SOAP response: %v", encodeErr)
	}
}

func main() {
	ctx := context.Background()
	serviceConfig := &config.ServiceConfig{}
	if configErr := envconfig.Process(ctx, serviceConfig); configErr != nil {
		log.Fatal(ctx, "failed to init config", configErr)
	}
	http.HandleFunc(serviceConfig.SoapGatewayConfig.Endpoint, soapHandler)

	address := fmt.Sprintf("%s:%s", serviceConfig.SoapGatewayConfig.EndpointHost, serviceConfig.SoapGatewayConfig.EndpointPort)
	log.Println(fmt.Sprintf("soap mock server running on address %s..", address))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", serviceConfig.SoapGatewayConfig.EndpointPort), nil))
}
