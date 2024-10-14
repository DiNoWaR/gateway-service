package gateway

import (
	"bytes"
	"encoding/xml"
	"fmt"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const callbackHeader = "Callback-URL"

type SoapGateway struct {
	Endpoint string
	Logger   *zap.Logger
}

func (sg *SoapGateway) ProcessDeposit(req DepositReq, callbackUrl string) (*DepositResponse, error) {
	soapReq, marshalErr := xml.MarshalIndent(Envelope{XMLName: xml.Name{}, Body: Body{DepositReq: &req}}, "", "  ")
	if marshalErr != nil {
		sg.Logger.Error("xml marshal failed", zap.Error(marshalErr))
		return &DepositResponse{}, marshalErr
	}

	resp, requestErr := makeRequest(sg.Endpoint, soapReq, callbackUrl, sg.Logger)
	if requestErr != nil {
		sg.Logger.Error("http request failed", zap.Error(requestErr))
		return &DepositResponse{}, requestErr
	}
	defer resp.Body.Close()

	responseBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		sg.Logger.Error("http response failed", zap.Error(readErr))
		return &DepositResponse{}, readErr
	}
	var envelope Envelope
	unmarshalErr := xml.Unmarshal(responseBytes, &envelope)
	if unmarshalErr != nil {
		sg.Logger.Error("xml unmarshal failed", zap.Error(unmarshalErr))
		return &DepositResponse{}, unmarshalErr
	}

	return envelope.Body.DepositResponse, nil
}

func (sg *SoapGateway) ProcessWithdrawal(req WithdrawReq, callbackUrl string) (*WithdrawResponse, error) {
	soapReq, marshalErr := xml.MarshalIndent(Envelope{XMLName: xml.Name{}, Body: Body{WithdrawReq: &req}}, "", "  ")
	if marshalErr != nil {
		sg.Logger.Error("xml marshal failed", zap.Error(marshalErr))
		return &WithdrawResponse{}, marshalErr
	}

	resp, requestErr := makeRequest(sg.Endpoint, soapReq, callbackUrl, sg.Logger)
	if requestErr != nil {
		sg.Logger.Error("http request failed", zap.Error(requestErr))
		return &WithdrawResponse{}, requestErr
	}
	defer resp.Body.Close()

	responseBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		sg.Logger.Error("http response failed", zap.Error(readErr))
		return &WithdrawResponse{}, readErr
	}
	var envelope Envelope
	unmarshalErr := xml.Unmarshal(responseBytes, &envelope)
	if unmarshalErr != nil {
		sg.Logger.Error("xml unmarshal failed", zap.Error(unmarshalErr))
		return &WithdrawResponse{}, unmarshalErr
	}

	return envelope.Body.WithdrawResponse, nil
}

func makeRequest(endpoint string, soapReq []byte, callbackUrl string, logger *zap.Logger) (*http.Response, error) {
	url := fmt.Sprintf("http://%s", endpoint)
	httpRequest, httpReqErr := http.NewRequest("POST", url, bytes.NewBuffer(soapReq))
	if httpReqErr != nil {
		logger.Error("failed to create HTTP request", zap.Error(httpReqErr))
		return nil, httpReqErr
	}

	httpRequest.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpRequest.Header.Set(callbackHeader, callbackUrl)

	client := &http.Client{}
	resp, requestErr := client.Do(httpRequest)

	if requestErr != nil {
		logger.Error("http request failed", zap.Error(requestErr))
		return resp, requestErr
	}
	return resp, nil
}
