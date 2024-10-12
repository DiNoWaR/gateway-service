package common

import (
	"github.com/shopspring/decimal"
	"time"
)

type Transaction struct {
	Id          string
	ReferenceId string
	AccountId   string
	Amount      decimal.Decimal
	Currency    string
	Status      TransactionStatus
	Operation   Operation
	Timestamp   time.Time
}

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "PENDING"
	StatusSuccess   TransactionStatus = "SUCCESS"
	StatusFailed    TransactionStatus = "FAILED"
	StatusCancelled TransactionStatus = "CANCELLED"
)

type Operation string

const (
	Withdraw Operation = "Withdraw"
	Deposit  Operation = "Deposit"
)
