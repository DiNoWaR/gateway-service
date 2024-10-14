package gateway

import . "github.com/dinowar/gateway-service/internal/pkg/domain/model"

type PaymentGateway interface {
	ProcessDeposit(req DepositReq, callbackUrl string) (*DepositResponse, error)
	ProcessWithdrawal(req WithdrawReq, callbackUrl string) (*WithdrawResponse, error)
}
