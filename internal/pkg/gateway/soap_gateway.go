package gateway

import (
	"bytes"
	"encoding/xml"
	"fmt"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	"github.com/dinowar/gateway-service/internal/pkg/util"
	"go.uber.org/zap"
	"io"
)

type SoapGateway struct {
	Endpoint        string
	Logger          *zap.Logger
	RetryInterval   int
	RetryElapseTime int
}

func (sg *SoapGateway) ProcessDeposit(req DepositReq, callbackUrl string) (*DepositResponse, error) {
	soapReq, marshalErr := xml.MarshalIndent(Envelope{XMLName: xml.Name{}, Body: Body{DepositReq: &req}}, "", "  ")
	if marshalErr != nil {
		sg.Logger.Error("xml marshal failed", zap.Error(marshalErr))
		return &DepositResponse{}, marshalErr
	}

	url := fmt.Sprintf("http://%s", sg.Endpoint)
	resp, retryErr := util.RetryableRequest(url, "POST", bytes.NewBuffer(soapReq), callbackUrl, "text/xml; charset=utf-8", sg.RetryInterval, sg.RetryElapseTime)
	if retryErr != nil {
		sg.Logger.Error("HandleDeposit: error processing deposit after retries: %v", zap.Error(retryErr))
		return &DepositResponse{}, retryErr
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

	url := fmt.Sprintf("http://%s", sg.Endpoint)
	resp, retryErr := util.RetryableRequest(url, "POST", bytes.NewBuffer(soapReq), callbackUrl, "text/xml; charset=utf-8", sg.RetryInterval, sg.RetryElapseTime)
	if retryErr != nil {
		sg.Logger.Error("HandleDeposit: error processing deposit after retries: %v", zap.Error(retryErr))
		return &WithdrawResponse{}, retryErr
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
