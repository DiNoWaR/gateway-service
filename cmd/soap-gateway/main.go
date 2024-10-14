package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/dinowar/gateway-service/internal/pkg/config"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	"github.com/google/uuid"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
)

var logger *zap.Logger

const (
	gatewayId = "soap"
	soapNS    = "http://schemas.xmlsoap.org/soap/envelope/"
	xsi       = "http://www.w3.org/2001/XMLSchema-instance"
	xsd       = "http://www.w3.org/2001/XMLSchema"
)

func init() {
	logger, _ = zap.NewDevelopment()
}

func soapHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		logger.Error("", zap.Error(readErr))
		http.Error(w, readErr.Error(), http.StatusBadRequest)
		return
	}

	var envelope Envelope
	unmarshalErr := xml.Unmarshal(bodyBytes, &envelope)
	if unmarshalErr != nil {
		logger.Error("", zap.Error(unmarshalErr))
		http.Error(w, unmarshalErr.Error(), http.StatusBadRequest)
		return
	}

	var response interface{}
	if envelope.Body.DepositReq != nil {
		req := envelope.Body.DepositReq
		fmt.Printf("Processing CashInRequest: %+v\n", req)

		response = DepositResponse{
			Gateway:       gatewayId,
			TransactionID: uuid.NewString(),
			Status:        StatusSuccess,
			Message:       "CashIn processed successfully",
			AccountID:     req.AccountID,
		}
	} else if envelope.Body.WithdrawReq != nil {
		req := envelope.Body.WithdrawReq
		fmt.Printf("Processing CashOutRequest: %+v\n", req)

		response = WithdrawResponse{
			Gateway:       gatewayId,
			TransactionID: uuid.NewString(),
			Status:        StatusSuccess,
			Message:       "CashOut processed successfully",
		}
	} else {
		http.Error(w, "Unknown request", http.StatusBadRequest)
		return
	}

	soapResponse := struct {
		XMLName xml.Name `xml:"soap:Envelope"`
		SoapNS  string   `xml:"xmlns:soap,attr"`
		XSI     string   `xml:"xmlns:xsi,attr"`
		XSD     string   `xml:"xmlns:xsd,attr"`
		Body    struct {
			XMLName xml.Name `xml:"soap:Body"`
			Content interface{}
		}
	}{
		SoapNS: soapNS,
		XSI:    xsi,
		XSD:    xsd,
	}
	soapResponse.Body.Content = response

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	xml.NewEncoder(w).Encode(soapResponse)
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
