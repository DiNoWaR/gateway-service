package service

import (
	"database/sql"
	"errors"
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
		`INSERT INTO transactions (reference_id, account_id, amount, currency, status, operation) 
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (reference_id) 
		 DO UPDATE SET account_id = EXCLUDED.account_id, amount = EXCLUDED.amount, currency = EXCLUDED.currency, 
		               status = EXCLUDED.status, operation = EXCLUDED.operation`,
		txn.ReferenceId, txn.AccountId, txn.Amount, txn.Currency, txn.Status, txn.Operation,
	)
	return trxSaveErr
}

func (rep *RepositoryService) GetTransaction(referenceId string) (Transaction, error) {
	var txn Transaction
	row := rep.db.QueryRow(`
		SELECT 
			COALESCE(id, '') AS id, 
			reference_id, 
			account_id, 
			amount, 
			currency, 
			status, 
			operation,
			ts 
		FROM transactions 
		WHERE reference_id = $1`, referenceId)

	trxErr := row.Scan(&txn.Id, &txn.ReferenceId, &txn.AccountId, &txn.Amount, &txn.Currency, &txn.Status, &txn.Operation, &txn.Ts)
	if errors.Is(trxErr, sql.ErrNoRows) {
		return Transaction{}, nil
	}
	if trxErr != nil {
		return Transaction{}, trxErr
	}

	return txn, nil
}
