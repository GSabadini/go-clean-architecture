package database

//NoSQLDbHandler expõe os métodos disponíveis para as abstrações de banco
type NoSQLDBHandler interface {
	Insert(string, interface{}) error
}
