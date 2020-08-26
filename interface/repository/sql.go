package repository

import "context"

type SQLHandler interface {
	ExecuteContext(context.Context, string, ...interface{}) error
	QueryContext(context.Context, string, ...interface{}) (Row, error)
	//BeginTx(ctx context.Context) (Tx, error)
}

type Row interface {
	Scan(dest ...interface{}) error
	Next() bool
	Err() error
	Close() error
}

type Tx interface {
	Commit() error
	Rollback() error
}
