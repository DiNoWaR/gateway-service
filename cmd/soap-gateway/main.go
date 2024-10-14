package main

import (
	"bytes"
	"context"
	"encoding/json"
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
	"time"
)

var (
	logger *zap.Logger
)

const (
	gatewayId      = "soap"
	soapNS         = "http://schemas.xmlsoap.org/soap/envelope/"
	xsi            = "http://www.w3.org/2001/XMLSchema-instance"
	xsd            = "http://www.w3.org/2001/XMLSchema"
	callbackHeader = "Callback-URL"
	delay          = 3 * time.Second
)

func init() {
	logger, _ = zap.NewDevelopment()
}

func soapHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		logger.Error("Error reading body", zap.Error(readErr))
		http.Error(w, readErr.Error(), http.StatusBadRequest)
		return
	}

	var envelope Envelope
	unmarshalErr := xml.Unmarshal(bodyBytes, &envelope)
	if unmarshalErr != nil {
		logger.Error("Error unmarshalling XML", zap.Error(unmarshalErr))
		http.Error(w, unmarshalErr.Error(), http.StatusBadRequest)
		return
	}

	var response interface{}
	if envelope.Body.DepositReq != nil {
		req := envelope.Body.DepositReq
		fmt.Printf("Processing CashInRequest: %+v\n", req)

		transactionId := uuid.NewString()
		response = DepositResponse{
			Gateway:       gatewayId,
			TransactionID: transactionId,
			Status:        StatusPending,
			Message:       "CashIn request received and is being processed",
			AccountID:     req.AccountID,
		}

		go processTransaction(transactionId, req.ReferenceID, r.Header.Get(callbackHeader))

	} else if envelope.Body.WithdrawReq != nil {
		req := envelope.Body.WithdrawReq
		fmt.Printf("Processing CashOutRequest: %+v\n", req)

		transactionId := uuid.NewString()
		response = WithdrawResponse{
			Gateway:       gatewayId,
			TransactionID: transactionId,
			Status:        StatusPending,
			Message:       "CashOut request received and is being processed",
		}

		go processTransaction(transactionId, req.ReferenceID, r.Header.Get(callbackHeader))

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

func processTransaction(transactionId, referenceId string, callbackURL string) {
	// imitation of delay
	time.Sleep(delay)
	message := fmt.Sprintf("%s processed successfully", transactionId)

	callbackData := map[string]string{
		"transaction_id": transactionId,
		"reference_id":   referenceId,
		"status":         "SUCCESS",
		"message":        message,
	}

	callbackErr := sendCallback(callbackURL, callbackData)
	if callbackErr != nil {
		logger.Error("error sending callback", zap.Error(callbackErr))
	} else {
		logger.Info("callback sent successfully")
	}
}

func sendCallback(callbackURL string, data map[string]string) error {
	if callbackURL == "" {
		logger.Error("callback url is empty")
		return fmt.Errorf("callback URL is empty")
	}

	logger.Info(fmt.Sprintf("callback url: %s", callbackURL))
	reqBody, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		logger.Error("error marshalling callback", zap.Error(marshalErr))
		return marshalErr
	}

	resp, respErr := http.Post(callbackURL, "application/json", bytes.NewBuffer(reqBody))
	if respErr != nil {
		logger.Error("error posting callback", zap.Error(respErr))
		return respErr
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("callback returned non-OK status: %s", resp.Status)
	}
	return nil
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
