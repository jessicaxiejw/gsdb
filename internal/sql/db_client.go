package sql

type DBClient interface {
	CreateTable(name string, columns []interface{}) error
	InsertRows(name string, values map[string][]interface{}) error
}
