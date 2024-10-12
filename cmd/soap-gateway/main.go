package main

import (
	"encoding/xml"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/common"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/soap_gateway"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

const (
	gatewayId = "soap_gateway"
	soapNS    = "http://schemas.xmlsoap.org/soap/envelope/"
	address   = "soap:9092"
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
	if envelope.Body.DepositReq != nil {
		req := envelope.Body.DepositReq
		log.Printf("processing deposit request: %+v\n\n", req)

		response = DepositResponse{
			Gateway:       gatewayId,
			TransactionID: uuid.NewString(),
			Status:        StatusSuccess,
			Message:       "deposit processed successfully",
			AccountID:     envelope.Body.DepositReq.AccountID,
		}
	} else if envelope.Body.WithdrawReq != nil {
		req := envelope.Body.WithdrawReq
		log.Printf("processing withdraw request: %+v\n", req)

		response = WithdrawResponse{
			Gateway:       gatewayId,
			TransactionID: uuid.NewString(),
			Status:        StatusSuccess,
			Message:       "withdraw processed successfully",
			AccountID:     envelope.Body.WithdrawReq.AccountID,
		}
	} else {
		http.Error(w, "unknown request", http.StatusBadRequest)
		return
	}

	soapResponse := struct {
		XMLName xml.Name `xml:"soap:Envelope"`
		SoapNS  string   `xml:"xmlns:soap,attr"`
		Body    struct {
			XMLName xml.Name `xml:"soap:Body"`
			Content interface{}
		}
	}{
		SoapNS: soapNS,
	}
	soapResponse.Body.Content = response

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	encodeErr := xml.NewEncoder(w).Encode(soapResponse)
	if encodeErr != nil {
		return
	}
}

func main() {
	http.HandleFunc("/soap", soapHandler)

	log.Println("soap mock server running on port 9002...")
	log.Fatal(http.ListenAndServe(":9002", nil))
}
