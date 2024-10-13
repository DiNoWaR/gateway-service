package gateway

type PaymentGateway interface {
	ProcessDeposit(req DepositReq) (DepositResponse, error)
	ProcessWithdrawal(req WithdrawReq) (WithdrawResponse, error)
}
