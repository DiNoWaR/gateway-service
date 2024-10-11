package util

import (
	"database/sql"
)

func InitDB() (*sql.DB, error) {
	db, dbErr := sql.Open("sqlite3", "./transactions.db")
	if dbErr != nil {
		return nil, dbErr
	}
	defer db.Close()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS transactions (
		id TEXT PRIMARY KEY,
		account_id TEXT,
		amount REAL,
		currency TEXT,
		status TEXT
	)`
	_, execErr := db.Exec(createTableQuery)
	if execErr != nil {
		return nil, execErr
	}
	return db, dbErr
}
