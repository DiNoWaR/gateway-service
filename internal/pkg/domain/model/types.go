package model

import (
	"encoding/xml"
	"github.com/shopspring/decimal"
	"time"
)

type Envelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	Body    Body     `xml:"soap:Body"`
}

type Body struct {
	Content interface{} `xml:",any"`
}

type ClientDepositRequest struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	AccountID string  `json:"account_id"`
	GatewayID string  `json:"gateway_id"`
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
