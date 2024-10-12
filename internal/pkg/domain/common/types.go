package common

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "PENDING"
	StatusSuccess   TransactionStatus = "SUCCESS"
	StatusFailed    TransactionStatus = "FAILED"
	StatusCancelled TransactionStatus = "CANCELLED"
)

type TransactionType string

const (
	Withdraw TransactionStatus = "Withdraw"
	Deposit  TransactionStatus = "Deposit"
)
