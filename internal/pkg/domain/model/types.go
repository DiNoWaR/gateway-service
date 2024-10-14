package model

import (
	"encoding/xml"
	"github.com/shopspring/decimal"
	"time"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
}

type Body struct {
	XMLName         xml.Name         `xml:"Body"`
	DepositReq      *DepositReq      `xml:"DepositRequest"`
	WithdrawReq     *WithdrawReq     `xml:"WithdrawRequest"`
	DepositResponse *DepositResponse `xml:"DepositResponse"`
}

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

type Transaction struct {
	Id          string
	ReferenceId string
	AccountId   string
	Amount      decimal.Decimal
	Currency    string
	Status      TransactionStatus
	Operation   Operation
	Ts          time.Time
}

type TransactionStatus string

const (
	StatusPending TransactionStatus = "PENDING"
	StatusSuccess TransactionStatus = "SUCCESS"
	StatusFailed  TransactionStatus = "FAILED"
)

type Operation string

const (
	Withdraw Operation = "Withdraw"
	Deposit  Operation = "Deposit"
)
