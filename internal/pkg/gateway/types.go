package gateway

import (
	"encoding/xml"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
)

type DepositReq struct {
	XMLName     xml.Name `xml:"DepositRequest"`
	Amount      float64  `json:"Amount" xml:"Amount"`
	Currency    string   `json:"Currency" xml:"Currency"`
	ReferenceID string   `json:"ReferenceId" xml:"ReferenceId"`
	AccountID   string   `json:"AccountId" xml:"AccountId"`
}

type DepositResponse struct {
	XMLName       xml.Name          `xml:"DepositResponse"`
	Gateway       string            `json:"Gateway" xml:"Gateway"`
	TransactionID string            `json:"TransactionId" xml:"TransactionId"`
	AccountID     string            `json:"AccountId" xml:"AccountId"`
	Status        TransactionStatus `json:"Status" xml:"Status"`
	Message       string            `json:"Message" xml:"Message"`
}

type WithdrawReq struct {
	XMLName     xml.Name `xml:"WithdrawRequest"`
	Amount      float64  `json:"Amount" xml:"Amount"`
	Currency    string   `json:"Currency" xml:"Currency"`
	ReferenceID string   `json:"ReferenceId" xml:"ReferenceId"`
	AccountID   string   `json:"AccountId" xml:"AccountId"`
}

type WithdrawResponse struct {
	XMLName       xml.Name          `xml:"WithdrawResponse"`
	Gateway       string            `json:"Gateway" xml:"Gateway"`
	TransactionID string            `json:"TransactionId" xml:"TransactionId"`
	AccountID     string            `json:"AccountId" xml:"AccountId"`
	Status        TransactionStatus `json:"Status" xml:"Status"`
	Message       string            `json:"Message" xml:"Message"`
}
