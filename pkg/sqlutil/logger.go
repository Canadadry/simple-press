package sqlutil

import (
	"app/pkg/middleware"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
}

type Logger struct {
	db DBTX
}

func NewLogger(db DBTX) *Logger {
	return &Logger{
		db: db,
	}
}

func (l *Logger) ExecContext(ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := l.db.ExecContext(ctx, query, params...)
	middleware.LogSQLQuery(ctx, "ExecContext", query, time.Since(start), params...)
	if err != nil {
		return nil, fmt.Errorf("while executing \"%s\" got %w", query, err)
	}
	return result, err
}

func (l *Logger) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	start := time.Now()
	stmt, err := l.db.PrepareContext(ctx, query)
	middleware.LogSQLQuery(ctx, "PrepareContext", query, time.Since(start))
	if err != nil {
		return nil, fmt.Errorf("while executing \"%s\" got %w", query, err)
	}
	return stmt, err
}

func (l *Logger) QueryContext(ctx context.Context, query string, params ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := l.db.QueryContext(ctx, query, params...)
	middleware.LogSQLQuery(ctx, "QueryContext", query, time.Since(start), params...)
	if err != nil {
		return nil, fmt.Errorf("while executing \"%s\" got %w", query, err)
	}
	return rows, err
}

func (l *Logger) QueryRowContext(ctx context.Context, query string, params ...interface{}) *sql.Row {
	start := time.Now()
	row := l.db.QueryRowContext(ctx, query, params...)
	middleware.LogSQLQuery(ctx, "QueryRowContext", query, time.Since(start), params...)
	return row
}

func (l *Logger) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return l.db.BeginTx(ctx, opts)
}
