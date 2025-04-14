package sqlutil

import (
	"context"
	"database/sql"
)

type TxWrapper struct {
	tx *sql.Tx
}

func (t TxWrapper) BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error) {
	return t.tx, nil
}

func (t TxWrapper) ExecContext(ctx context.Context, s string, i ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, s, i...)
}

func (t TxWrapper) PrepareContext(ctx context.Context, s string) (*sql.Stmt, error) {
	return t.tx.PrepareContext(ctx, s)
}

func (t TxWrapper) QueryContext(ctx context.Context, s string, i ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, s, i...)
}

func (t TxWrapper) QueryRowContext(ctx context.Context, s string, i ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, s, i...)
}

func TxWrap(tx *sql.Tx) DBTX {
	return TxWrapper{
		tx: tx,
	}
}
