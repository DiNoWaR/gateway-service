package gateway

import . "github.com/dinowar/gateway-service/internal/pkg/domain/model"

type PaymentGateway interface {
	ProcessDeposit(req DepositReq) (*DepositResponse, error)
	ProcessWithdrawal(req WithdrawReq) (*WithdrawResponse, error)
}
