package sql

import (
	"fmt"
	"reflect"

	pg_query "github.com/pganalyze/pg_query_go/v4"
)

type PostgreSQL struct {
	client DBClient
}

func NewPostgreSQL(client DBClient) *PostgreSQL {
	return &PostgreSQL{client: client}
}

func (p *PostgreSQL) Execute(statement string) error {
	result, err := pg_query.Parse(statement)
	if err != nil {
		return err // TODO: wrap this error
	}
	fmt.Println(pg_query.ParseToJSON(statement)) // TODO: delete
	for _, stmt := range result.GetStmts() {
		node := stmt.GetStmt().GetNode()
		switch node.(type) {
		case *pg_query.Node_CreateStmt:
			createStmt := node.(*pg_query.Node_CreateStmt).CreateStmt
			err = p.createTable(createStmt)
		case *pg_query.Node_InsertStmt:
			insertStmt := node.(*pg_query.Node_InsertStmt).InsertStmt
			err = p.insert(insertStmt)
		default:
			err = fmt.Errorf("unfortunately, we do not support %s at this time", reflect.TypeOf(node))
		}
	}
	return err
}

func (p *PostgreSQL) createTable(stmt *pg_query.CreateStmt) error {
	tableName := stmt.GetRelation().GetRelname()

	elts := stmt.GetTableElts()
	columns := make([]interface{}, len(elts))
	for index, elt := range elts {
		columnDef := elt.GetColumnDef()
		columns[index] = columnDef.GetColname()
	}

	return p.client.CreateTable(tableName, columns)
}

func (p *PostgreSQL) insert(stmt *pg_query.InsertStmt) error {
	tableName := stmt.GetRelation().GetRelname()

	cols := stmt.GetCols()
	columnNames := make([]string, len(cols))
	for index, col := range cols {
		columnNames[index] = col.GetResTarget().Name
	}

	valueLists := stmt.GetSelectStmt().GetSelectStmt().GetValuesLists()
	values := make([][]interface{}, len(valueLists))
	for i, valueList := range valueLists {
		items := valueList.GetList().GetItems()
		values[i] = make([]interface{}, len(items))
		for j, item := range items {
			values[i][j] = item.GetAConst().GetSval().GetSval()
		}
	}

	return p.client.InsertRows(tableName, columnNames, values)
}
