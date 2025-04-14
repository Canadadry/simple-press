package sqlutil

import (
	"app/pkg/dbconn"
	"context"
	"database/sql"
)

type AutoClose struct {
	Filename string
}

func NewAutoCloseConn(filename string, autoclose bool) (dbconn.DBTX, error) {
	if !autoclose {
		return dbconn.Open(filename)
	}
	return &AutoClose{
		Filename: filename,
	}, nil
}

func (a *AutoClose) ExecContext(ctx context.Context, q string, args ...any) (sql.Result, error) {
	dbtx, err := dbconn.Open(a.Filename)
	if err != nil {
		return nil, err
	}
	defer dbtx.Close()
	return dbtx.ExecContext(ctx, q, args...)
}

func (a *AutoClose) PrepareContext(ctx context.Context, q string) (*sql.Stmt, error) {
	dbtx, err := dbconn.Open(a.Filename)
	if err != nil {
		return nil, err
	}
	defer dbtx.Close()
	return dbtx.PrepareContext(ctx, q)
}

func (a *AutoClose) QueryContext(ctx context.Context, q string, args ...any) (*sql.Rows, error) {
	dbtx, err := dbconn.Open(a.Filename)
	if err != nil {
		return nil, err
	}
	defer dbtx.Close()
	return dbtx.QueryContext(ctx, q, args...)
}

func (a *AutoClose) QueryRowContext(ctx context.Context, q string, args ...any) *sql.Row {
	dbtx, err := dbconn.Open(a.Filename)
	if err != nil {
		return nil
	}
	defer dbtx.Close()
	return dbtx.QueryRowContext(ctx, q, args...)
}

func (a *AutoClose) BeginTx(ctx context.Context, opt *sql.TxOptions) (*sql.Tx, error) {
	dbtx, err := dbconn.Open(a.Filename)
	if err != nil {
		return nil, err
	}
	defer dbtx.Close()
	return dbtx.BeginTx(ctx, opt)
}

func (a *AutoClose) Ping() error {
	dbtx, err := dbconn.Open(a.Filename)
	if err != nil {
		return err
	}
	defer dbtx.Close()
	return dbtx.Ping()
}

func (a *AutoClose) Close() error {
	return nil
}
