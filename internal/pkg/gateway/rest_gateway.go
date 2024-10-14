package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type RestGateway struct {
	BaseURL string
	Logger  *zap.Logger
}

func (rg *RestGateway) ProcessDeposit(req DepositReq, callbackUrl string) (*DepositResponse, error) {
	url := fmt.Sprintf("http://%s/deposit", rg.BaseURL)
	jsonData, marshalErr := json.Marshal(req)
	if marshalErr != nil {
		rg.Logger.Error(marshalErr.Error())
		return &DepositResponse{}, marshalErr
	}

	httpReq, httpReqErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if httpReqErr != nil {
		rg.Logger.Error("failed to create request", zap.String("url", url), zap.Error(httpReqErr))
		return &DepositResponse{}, httpReqErr
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set(callbackHeader, callbackUrl)

	client := &http.Client{}
	resp, requestErr := client.Do(httpReq)
	if requestErr != nil {
		rg.Logger.Error("http request failed", zap.String("url", url), zap.Error(requestErr))
		return &DepositResponse{}, requestErr
	}
	defer resp.Body.Close()

	responseBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		rg.Logger.Error("http response failed", zap.Error(readErr))
		return &DepositResponse{}, readErr
	}

	var depositResp DepositResponse
	decodeErr := json.Unmarshal(responseBytes, &depositResp)
	if decodeErr != nil {
		rg.Logger.Error(decodeErr.Error(), zap.String("url", url))
		return &DepositResponse{}, decodeErr
	}

	return &depositResp, nil
}

func (rg *RestGateway) ProcessWithdrawal(req WithdrawReq, callbackUrl string) (*WithdrawResponse, error) {
	url := fmt.Sprintf("http://%s/withdraw", rg.BaseURL)
	jsonData, marshalErr := json.Marshal(req)
	if marshalErr != nil {
		rg.Logger.Error(marshalErr.Error())
		return &WithdrawResponse{}, marshalErr
	}

	httpReq, httpReqErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if httpReqErr != nil {
		rg.Logger.Error("failed to create request", zap.String("url", url), zap.Error(httpReqErr))
		return &WithdrawResponse{}, httpReqErr
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set(callbackHeader, callbackUrl)

	client := &http.Client{}
	resp, requestErr := client.Do(httpReq)
	if requestErr != nil {
		rg.Logger.Error("http request failed", zap.String("url", url), zap.Error(requestErr))
		return &WithdrawResponse{}, requestErr
	}
	defer resp.Body.Close()

	responseBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		rg.Logger.Error("http response failed", zap.Error(readErr))
		return &WithdrawResponse{}, readErr
	}

	var withdrawResp WithdrawResponse
	decodeErr := json.Unmarshal(responseBytes, &withdrawResp)
	if decodeErr != nil {
		rg.Logger.Error(decodeErr.Error(), zap.String("url", url))
		return &WithdrawResponse{}, decodeErr
	}

	return &withdrawResp, nil
}
