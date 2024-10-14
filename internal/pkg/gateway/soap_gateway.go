package gateway

import (
	"bytes"
	"encoding/xml"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type SoapGateway struct {
	Endpoint string
	Logger   *zap.Logger
}

func (sg *SoapGateway) ProcessDeposit(req DepositReq) (*DepositResponse, error) {
	soapReq, marshalErr := xml.MarshalIndent(Envelope{
		XMLName: xml.Name{},
		Body: Body{
			DepositReq: &req,
		},
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
	}, "", "  ")
	if marshalErr != nil {
		sg.Logger.Error("xml marshal failed", zap.Error(marshalErr))
		return &DepositResponse{}, marshalErr
	}

	resp, requestErr := http.Post(sg.Endpoint, "text/xml; charset=utf-8", bytes.NewBuffer(soapReq))
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

func (sg *SoapGateway) ProcessWithdrawal(req WithdrawReq) (*WithdrawResponse, error) {
	return &WithdrawResponse{}, nil
}
