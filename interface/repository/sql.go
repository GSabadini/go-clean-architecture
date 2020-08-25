package repository

import "context"

//SQLHandler expõe os métodos disponíveis para as abstrações de banco SQL
type SQLHandler interface {
	ExecuteContext(context.Context, string, ...interface{}) error
	QueryContext(context.Context, string, ...interface{}) (Row, error)
	//BeginTx(ctx context.Context) (Tx, error)
}

//Row expõe os métodos disponíveis para as abstrações de linhas de banco SQL
type Row interface {
	Scan(dest ...interface{}) error
	Next() bool
	Err() error
	Close() error
}

//Tx expõe os métodos disponíveis para as abstrações de transações
type Tx interface {
	Commit() error
	Rollback() error
}
