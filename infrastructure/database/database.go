package database

//NoSQLDbHandler expõe os métodos disponíveis para as abstrações de banco
type NoSQLDbHandler interface {
	Store(string, interface{}) error
	Update(string, interface{}, interface{}) error
	FindAll(string, interface{}, interface{}) error
	FindOne(string, interface{}, interface{}, interface{}) error
}

//DbHandler expõe os métodos disponíveis para as abstrações de banco
type DbHandler interface {
	Store(string, ...interface{}) error
	Update(string, ...interface{}) error
	FindAll(string, ...interface{}) error
	FindOne(string, ...interface{}) error
}
