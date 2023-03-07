package sql

type DBClient interface {
	CreateTable(name string, columns []TableColumn) error
}

type TableColumn struct {
	name     string
	dataType string
}
