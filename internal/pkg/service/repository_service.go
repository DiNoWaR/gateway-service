package service

import (
	"database/sql"
	. "github.com/dinowar/gateway-service/internal/pkg/domain"
	_ "github.com/mattn/go-sqlite3"
)

type RepositoryService struct {
	db *sql.DB
}

func NewRepositoryService(db *sql.DB) *RepositoryService {
	return &RepositoryService{db: db}
}

func (rep *RepositoryService) SaveTransaction(txn Transaction) error {
	_, trxSaveErr := rep.db.Exec(
		"INSERT OR REPLACE INTO transactions (id, account_id, amount, currency, status) VALUES (?, ?, ?, ?)",
		txn.ID, txn.AccountId, txn.Amount, txn.Currency, txn.Status,
	)
	return trxSaveErr
}

func (rep *RepositoryService) GetTransaction(transactionID string) (Transaction, error) {
	var txn Transaction
	row := rep.db.QueryRow("SELECT id, account_id, amount, currency, status, operation FROM transactions WHERE id = ?", transactionID)
	trxErr := row.Scan(&txn.ID, &txn.AccountId, &txn.Amount, &txn.Currency, &txn.Status, &txn.Operation)
	if trxErr != nil {
		return Transaction{}, trxErr
	}
	return txn, nil
}

func (rep *RepositoryService) GetTransactions(accountId string) ([]Transaction, error) {
	var transactions []Transaction
	rows, trxErr := rep.db.Query("SELECT id, account_id, amount, currency, status, operation FROM transactions where account_id = ?", accountId)
	if trxErr != nil {
		return transactions, trxErr
	}
	defer rows.Close()

	for rows.Next() {
		var txn Transaction
		scanErr := rows.Scan(&txn.ID, &txn.AccountId, &txn.Amount, &txn.Currency, &txn.Status, &txn.Operation)
		if scanErr != nil {
			return transactions, scanErr
		}
		transactions = append(transactions, txn)
	}
	return transactions, nil
}
