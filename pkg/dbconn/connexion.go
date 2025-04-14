package dbconn

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const enableForeignKeys = "PRAGMA foreign_keys = ON"

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	Ping() error
	Close() error
}

func MustOpen(file string) DBTX {
	db, err := Open(file)
	if err != nil {
		panic(fmt.Sprintf("cannot init dbconn : %s", file))
	}
	return db
}

func Open(file string) (DBTX, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("cannot open database : %w", err)
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("while pinging database %w", err)
	}
	_, err = db.ExecContext(context.Background(), enableForeignKeys)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("while enabling fk check %w", err)
	}
	return db, nil
}
