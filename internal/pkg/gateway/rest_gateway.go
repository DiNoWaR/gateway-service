package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type RestGateway struct {
	BaseURL string
}

func (rg *RestGateway) ProcessDeposit(req DepositReq) (DepositResponse, error) {
	url := fmt.Sprintf("%s/deposit", rg.BaseURL)
	jsonData, marshalErr := json.Marshal(req)
	if marshalErr != nil {
		return DepositResponse{}, marshalErr
	}

	resp, requestErr := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if requestErr != nil {
		return DepositResponse{}, requestErr
	}
	defer resp.Body.Close()

	var depositResp DepositResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&depositResp)
	if decodeErr != nil {
		return DepositResponse{}, decodeErr
	}

	return depositResp, nil
}

func (rg *RestGateway) ProcessWithdrawal(req WithdrawReq) (WithdrawResponse, error) {
	url := fmt.Sprintf("%s/withdraw", rg.BaseURL)
	jsonData, marshalErr := json.Marshal(req)
	if marshalErr != nil {
		return WithdrawResponse{}, marshalErr
	}

	resp, requestErr := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if requestErr != nil {
		return WithdrawResponse{}, requestErr
	}
	defer resp.Body.Close()

	var withdrawResp WithdrawResponse
	decodeErr := json.NewDecoder(resp.Body).Decode(&withdrawResp)
	if decodeErr != nil {
		return WithdrawResponse{}, decodeErr
	}

	return withdrawResp, nil
}
