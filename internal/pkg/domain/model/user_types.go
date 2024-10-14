package model

type ClientDepositRequest struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	AccountID string  `json:"account_id"`
	GatewayID string  `json:"gateway_id"`
}
