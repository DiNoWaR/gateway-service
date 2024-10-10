package domain

type DepositRequest struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Gateway   string  `json:"gateway"`
}

type WithdrawRequest struct {
	AccountID string  `json:"account_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Gateway   string  `json:"gateway"`
}

type CallbackResponse struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
}

type Transaction struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Status   string  `json:"status"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
