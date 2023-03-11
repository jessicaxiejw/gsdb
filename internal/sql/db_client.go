package sql

type DBClient interface {
	CreateTable(name string, columns []interface{}) error
	InsertRows(name string, columnNames []string, values [][]interface{}) error
}
