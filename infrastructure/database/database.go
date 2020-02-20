package database

//NoSQLDbHandler expõe os métodos disponíveis para as abstrações de banco
type NoSQLDBHandler interface {
	Store(string, interface{}) error
	FindAll(string, interface{}, interface{}) error
}
