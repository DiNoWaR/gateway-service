package service

import (
	. "github.com/dinowar/gateway-service/internal/pkg/domain"
	"log"
)

type LogService struct{}

func NewLogService() *LogService {
	return &LogService{}
}

func (logger *LogService) LogTransaction(txn Transaction) {
	log.Printf("Transaction logged: ID=%s, Amount=%.2f, Currency=%s, Status=%s", txn.ID, txn.Amount, txn.Currency, txn.Status)
}
