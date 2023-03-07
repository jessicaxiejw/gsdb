package sql

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v4"
)

type PostgreSQL struct {
	client DBClient
}

func (p *PostgreSQL) Parse(statement string) error {
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
			p.createTable(createStmt)
		default:
			fmt.Println("do nothing") // TODO: any error for unknown type
		}
	}
	return nil
}

func (p *PostgreSQL) createTable(stmt *pg_query.CreateStmt) error {
	tableName := stmt.GetRelation().GetRelname()

	elts := stmt.GetTableElts()
	columns := make([]TableColumn, len(elts))
	for index, elt := range elts {
		columnDef := elt.GetColumnDef()
		columns[index] = TableColumn{
			name:     columnDef.GetColname(),
			dataType: columnDef.GetTypeName().GetNames()[1].String(),
		}
	}

	return p.client.CreateTable(tableName, columns)
}
