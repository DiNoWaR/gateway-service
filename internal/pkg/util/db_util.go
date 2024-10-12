package util

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func InitDB(host, port, dbname, user, password string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, dbErr := sql.Open("postgres", connStr)
	if dbErr != nil {
		return nil, dbErr
	}
	return db, nil
}
