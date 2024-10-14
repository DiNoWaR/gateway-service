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
		`INSERT INTO transactions (reference_id, account_id, amount, currency, status, operation, gateway_id) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (reference_id) 
		 DO UPDATE SET account_id = EXCLUDED.account_id, amount = EXCLUDED.amount, currency = EXCLUDED.currency, 
		               status = EXCLUDED.status, operation = EXCLUDED.operation`,
		txn.ReferenceId, txn.AccountId, txn.Amount, txn.Currency, txn.Status, txn.Operation, txn.GatewayId,
	)
	return trxSaveErr
}

func (rep *RepositoryService) UpdateTransaction(txn *Transaction) error {
	_, trxSaveErr := rep.db.Exec(
		`INSERT INTO transactions (id, reference_id, account_id, amount, currency, status, operation, message, gateway_id) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 ON CONFLICT (reference_id) 
		 DO UPDATE SET id = EXCLUDED.id, 
		               status = EXCLUDED.status, 
		               message = EXCLUDED.message`,
		txn.Id, txn.ReferenceId, txn.AccountId, txn.Amount, txn.Currency, txn.Status, txn.Operation, txn.Message, txn.GatewayId,
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
			COALESCE(message, '') AS message, 
			gateway_id,
			ts 
		FROM transactions 
		WHERE reference_id = $1`, referenceId)

	trxErr := row.Scan(&txn.Id, &txn.ReferenceId, &txn.AccountId, &txn.Amount, &txn.Currency, &txn.Status, &txn.Operation, &txn.Message, &txn.GatewayId, &txn.Ts)
	if errors.Is(trxErr, sql.ErrNoRows) {
		return Transaction{}, nil
	}
	if trxErr != nil {
		return Transaction{}, trxErr
	}

	return txn, nil
}

func (rep *RepositoryService) GetTransactions(accountId string) ([]Transaction, error) {
	rows, rowsErr := rep.db.Query(`
		SELECT 
			COALESCE(id, '') AS id,  
			reference_id, 
			account_id, 
			amount, 
			currency, 
			status, 
			operation,
			COALESCE(message, '') AS message, 
			gateway_id,
			ts 
		FROM transactions 
		WHERE account_id = $1 ORDER BY ts DESC`, accountId)
	if rowsErr != nil {
		return nil, rowsErr
	}

	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var txn Transaction
		cursorErr := rows.Scan(&txn.Id, &txn.ReferenceId, &txn.AccountId, &txn.Amount, &txn.Currency, &txn.Status, &txn.Operation, &txn.Message, &txn.GatewayId, &txn.Ts)
		if cursorErr != nil {
			return nil, cursorErr
		}
		transactions = append(transactions, txn)
	}

	if execErr := rows.Err(); execErr != nil {
		return nil, execErr
	}

	return transactions, nil
}
