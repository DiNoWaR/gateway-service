package service

import (
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dinowar/gateway-service/internal/pkg/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestSaveTransaction_Success(t *testing.T) {
	db, mock, mockErr := sqlmock.New()
	assert.NoError(t, mockErr)
	defer db.Close()

	rep := NewRepositoryService(db)

	txn := &model.Transaction{
		ReferenceId: "ref123",
		AccountId:   "ACC123",
		Amount:      decimal.NewFromFloat(100.50),
		Currency:    "USD",
		Status:      model.StatusPending,
		Operation:   model.Deposit,
		GatewayId:   "rest_gateway",
	}

	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(txn.ReferenceId, txn.AccountId, txn.Amount, txn.Currency, txn.Status, txn.Operation, txn.GatewayId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	saveErr := rep.SaveTransaction(txn)
	assert.NoError(t, saveErr)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveTransaction_Failure(t *testing.T) {
	db, mock, mockErr := sqlmock.New()
	assert.NoError(t, mockErr)
	defer db.Close()

	rep := NewRepositoryService(db)
	txn := &model.Transaction{
		ReferenceId: "ref123",
		AccountId:   "ACC123",
		Amount:      decimal.NewFromFloat(100.50),
		Currency:    "USD",
		Status:      model.StatusPending,
		Operation:   model.Deposit,
		GatewayId:   "rest_gateway",
	}

	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(txn.ReferenceId, txn.AccountId, txn.Amount, txn.Currency, txn.Status, txn.Operation, txn.GatewayId).
		WillReturnError(errors.New("failed to insert transaction"))

	saveErr := rep.SaveTransaction(txn)
	assert.Error(t, saveErr)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTransaction_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rep := NewRepositoryService(db)

	txn := &model.Transaction{
		Id:          "1",
		ReferenceId: "ref123",
		AccountId:   "ACC123",
		Amount:      decimal.NewFromFloat(100.50),
		Currency:    "USD",
		Status:      model.StatusSuccess,
		Operation:   model.Deposit,
		Message:     "Transaction successful",
		GatewayId:   "rest_gateway",
		Ts:          time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "reference_id", "account_id", "amount", "currency", "status", "operation", "message", "gateway_id", "ts"}).
		AddRow(txn.Id, txn.ReferenceId, txn.AccountId, txn.Amount, txn.Currency, txn.Status, txn.Operation, txn.Message, txn.GatewayId, txn.Ts)

	mock.ExpectQuery(`SELECT (.+) FROM transactions WHERE reference_id = ?`).
		WithArgs("ref123").
		WillReturnRows(rows)

	result, err := rep.GetTransaction("ref123")
	assert.NoError(t, err)
	assert.Equal(t, txn.ReferenceId, result.ReferenceId)
	assert.Equal(t, txn.Status, result.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTransaction_NotFound(t *testing.T) {
	db, mock, mockErr := sqlmock.New()
	assert.NoError(t, mockErr)
	defer db.Close()

	rep := NewRepositoryService(db)

	mock.ExpectQuery(`SELECT (.+) FROM transactions WHERE reference_id = ?`).
		WithArgs("ref123").
		WillReturnError(sql.ErrNoRows)

	_, txErr := rep.GetTransaction("ref123")
	assert.NoError(t, txErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}
