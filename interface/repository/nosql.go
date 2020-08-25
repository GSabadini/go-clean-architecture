package repository

import "context"

//NoSQLHandler expõe os métodos disponíveis para as abstrações de banco NoSQL
type NoSQLHandler interface {
	Store(context.Context, string, interface{}) error
	Update(context.Context, string, interface{}, interface{}) error
	FindAll(context.Context, string, interface{}, interface{}) error
	FindOne(context.Context, string, interface{}, interface{}, interface{}) error
}
