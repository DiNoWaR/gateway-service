package rest_gateway

import . "github.com/dinowar/gateway-service/internal/pkg/domain/common"

type DepositReq struct {
	Amount      float64 `json:"Amount"`
	Currency    string  `json:"Currency"`
	ReferenceID string  `json:"ReferenceId"`
	AccountID   string  `json:"AccountId"`
}

type WithdrawReq struct {
	Amount      float64 `json:"Amount"`
	Currency    string  `json:"Currency"`
	ReferenceID string  `json:"ReferenceId"`
	AccountID   string  `json:"AccountId"`
}

type DepositResponse struct {
	Gateway       string            `json:"Gateway"`
	TransactionID string            `json:"TransactionId"`
	AccountID     string            `json:"AccountId"`
	Status        TransactionStatus `json:"Status"`
	Message       string            `json:"Message"`
}

type WithdrawResponse struct {
	Gateway       string            `json:"Gateway"`
	TransactionID string            `json:"TransactionId"`
	AccountID     string            `json:"AccountId"`
	Status        TransactionStatus `json:"Status"`
	Message       string            `json:"Message"`
}
