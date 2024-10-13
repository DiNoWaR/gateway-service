package service

import (
	"database/sql"
	. "github.com/dinowar/gateway-service/internal/pkg/domain/model"
)

type RepositoryService struct {
	db *sql.DB
}

func NewRepositoryService(db *sql.DB) *RepositoryService {
	return &RepositoryService{db: db}
}

func (rep *RepositoryService) SaveTransaction(txn *Transaction) error {
	_, trxSaveErr := rep.db.Exec(
		"INSERT OR REPLACE INTO transactions (id, reference_id, account_id, amount, currency, status, operation) VALUES (?, ?, ?, ?, ?, ?, ?)",
		txn.Id, txn.AccountId, txn.Amount, txn.Currency, txn.Status,
	)
	return trxSaveErr
}

func (rep *RepositoryService) GetTransaction(referenceId string) (Transaction, error) {
	var txn Transaction
	row := rep.db.QueryRow("SELECT id, reference_id, account_id, amount, currency, status, operation FROM transactions WHERE reference_id = ?", referenceId)
	trxErr := row.Scan(&txn.Id, &txn.ReferenceId, &txn.AccountId, &txn.Amount, &txn.Currency, &txn.Status, &txn.Operation)
	if trxErr != nil {
		return Transaction{}, trxErr
	}
	return txn, nil
}
