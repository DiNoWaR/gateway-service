package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
	"github.com/dinowar/gateway-service/internal/pkg/util"
	"go.uber.org/zap"
	"io"
)

type RestGateway struct {
	BaseURL         string
	Logger          *zap.Logger
	RetryInterval   int
	RetryElapseTime int
}

func (rg *RestGateway) ProcessDeposit(req DepositReq, callbackUrl string) (*DepositResponse, error) {
	url := fmt.Sprintf("http://%s/deposit", rg.BaseURL)
	jsonData, marshalErr := json.Marshal(req)
	if marshalErr != nil {
		rg.Logger.Error(marshalErr.Error())
		return &DepositResponse{}, marshalErr
	}

	resp, retryErr := util.RetryableRequest(url, "POST", bytes.NewBuffer(jsonData), callbackUrl, "application/json", rg.RetryInterval, rg.RetryElapseTime)
	if retryErr != nil {
		rg.Logger.Error("HandleDeposit: error processing deposit after retries: %v", zap.Error(retryErr))
		return &DepositResponse{}, retryErr
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

	resp, retryErr := util.RetryableRequest(url, "POST", bytes.NewBuffer(jsonData), callbackUrl, "application/json", 1, 1)
	if retryErr != nil {
		rg.Logger.Error("HandleDeposit: error processing deposit after retries: %v", zap.Error(retryErr))
		return &WithdrawResponse{}, retryErr
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
