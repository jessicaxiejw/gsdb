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

func (p *PostgreSQL) Execute(statement string) (string, error) {
	result, err := pg_query.Parse(statement)
	if err != nil {
		return "", err // TODO: wrap this error
	}
	fmt.Println(pg_query.ParseToJSON(statement)) // TODO: delete
	output := ""
	for _, stmt := range result.GetStmts() {
		node := stmt.GetStmt().GetNode()
		switch node.(type) {
		case *pg_query.Node_CreateStmt:
			createStmt := node.(*pg_query.Node_CreateStmt).CreateStmt
			err = p.createTable(createStmt)
		case *pg_query.Node_InsertStmt:
			insertStmt := node.(*pg_query.Node_InsertStmt).InsertStmt
			err = p.insert(insertStmt)
		case *pg_query.Node_SelectStmt:
			selectStmt := node.(*pg_query.Node_SelectStmt).SelectStmt
			output, err = p.query(selectStmt)
		default:
			err = fmt.Errorf("unfortunately, we do not support %s at this time", reflect.TypeOf(node))
		}
	}
	return output, err
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

	values := map[string][]interface{}{} // key is the column name, the value is the array of values under that column
	valueLists := stmt.GetSelectStmt().GetSelectStmt().GetValuesLists()
	indexToColumnMapping := map[int]string{}
	for index, col := range stmt.GetCols() {
		name := col.GetResTarget().Name
		indexToColumnMapping[index] = name
		values[name] = make([]interface{}, len(valueLists))
	}

	for i, valueList := range valueLists {
		for j, item := range valueList.GetList().GetItems() {
			colName := indexToColumnMapping[j]
			values[colName][i] = item.GetAConst().GetSval().GetSval()
		}
	}

	return p.client.InsertRows(tableName, values)
}

func (p *PostgreSQL) query(stmt *pg_query.SelectStmt) (string, error) {
	// tableName := stmt.GetFromClause()[0].GetRangeVar().GetRelname() // TODO: make JOIN, UNION work

	// columns := make([]string, len(stmt.GetTargetList()))
	// for index, target := range stmt.GetTargetList() {
	// 	columns[index] = target.GetResTarget().GetName()
	// }

	// where := map[string][]interface{}{}
	// whereClause := stmt.GetWhereClause()
	// if whereClause != nil {
	// }
	return "", nil
}
