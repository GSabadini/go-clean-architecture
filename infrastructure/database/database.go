package database

//NoSQLHandler expõe os métodos disponíveis para as abstrações de banco NoSQL
type NoSQLHandler interface {
	Store(string, interface{}) error
	Update(string, interface{}, interface{}) error
	FindAll(string, interface{}, interface{}) error
	FindOne(string, interface{}, interface{}, interface{}) error
}

//SQLHandler expõe os métodos disponíveis para as abstrações de banco SQL
type SQLHandler interface {
	Execute(string, ...interface{}) error
	Query(string, ...interface{}) (Row, error)
}

//Row
type Row interface {
	Scan(dest ...interface{}) error
	Next() bool
	Err() error
}
