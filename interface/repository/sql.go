package repository

import "context"

type SQLHandler interface {
	ExecuteContext(context.Context, string, ...interface{}) error
	QueryContext(context.Context, string, ...interface{}) (Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) Row
	//BeginTx(ctx context.Context) (Tx, error)
}

type Rows interface {
	Scan(dest ...interface{}) error
	Next() bool
	Err() error
	Close() error
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Tx interface {
	Commit() error
	Rollback() error
}
