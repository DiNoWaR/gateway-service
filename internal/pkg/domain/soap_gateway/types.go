package soap_gateway

import (
	"encoding/xml"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/common"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	XMLName     xml.Name     `xml:"Body"`
	DepositReq  *DepositReq  `xml:"DepositRequest"`
	WithdrawReq *WithdrawReq `xml:"WithdrawRequest"`
}

type DepositReq struct {
	Amount      float64 `xml:"Amount"`
	Currency    string  `xml:"Currency"`
	ReferenceID string  `xml:"ReferenceId"`
	AccountID   string  `xml:"AccountId"`
}

type WithdrawReq struct {
	Amount      float64 `xml:"Amount"`
	Currency    string  `xml:"Currency"`
	ReferenceID string  `xml:"ReferenceId"`
	AccountID   string  `xml:"AccountId"`
}

type DepositResponse struct {
	XMLName       xml.Name          `xml:"DepositResponse"`
	Gateway       string            `xml:"Gateway"`
	TransactionID string            `xml:"TransactionId"`
	AccountID     string            `xml:"AccountId"`
	Status        TransactionStatus `xml:"Status"`
	Message       string            `xml:"Message"`
}

type WithdrawResponse struct {
	XMLName       xml.Name          `xml:"WithdrawResponse"`
	TransactionID string            `xml:"TransactionId"`
	Gateway       string            `xml:"Gateway"`
	AccountID     string            `xml:"AccountId"`
	Status        TransactionStatus `xml:"Status"`
	Message       string            `xml:"Message"`
}
