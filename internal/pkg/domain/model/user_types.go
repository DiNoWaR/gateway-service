package model

type ClientRequest struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	AccountID string  `json:"account_id"`
	GatewayID string  `json:"gateway_id"`
}

type GetTransactionRequest struct {
	ReferenceId string `json:"reference_id"`
}

type GetTransactionsRequest struct {
	AccountId string `json:"account_id"`
}

type CallbackPayload struct {
	TransactionId string `json:"transaction_id"`
	ReferenceId   string `json:"reference_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}
