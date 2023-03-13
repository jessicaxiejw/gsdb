package sheet

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/samber/lo"
)

type table struct {
	columnToIndex map[string]int // key is the column name, value is the index
	content       [][]string     // first index is the row, second index is the column
	errs          error

	selectedRowIndices []int
	selectedColumns    []string
}

func newTableFromGoogleSheet(values [][]interface{}) *table {
	columnToIndex := map[string]int{}
	for i, name := range values[0] {
		columnToIndex[name.(string)] = i
	}

	content := make([][]string, len(values)-1)
	selectedRowIndices := make([]int, len(values)-1)
	for i, row := range values[1:] {
		selectedRowIndices[i] = i
		content[i] = make([]string, len(row))
		for j, value := range row {
			content[i][j] = fmt.Sprint(value)
		}
	}

	return &table{
		columnToIndex:      columnToIndex,
		content:            content,
		selectedRowIndices: selectedRowIndices,
		selectedColumns:    []string{},
	}
}

func newTable() *table {
	return &table{
		columnToIndex:      map[string]int{},
		content:            [][]string{},
		selectedRowIndices: []int{},
		selectedColumns:    []string{},
	}
}

func (t *table) Select(columns []string) *table {
	t.selectedColumns = columns
	return t
}

func (t *table) Where(column, operator, value string) *table {
	newSelectedRowIndices := []int{}
	columnIndex := t.columnToIndex[column]
	for _, rowIndex := range t.selectedRowIndices {
		row := t.content[rowIndex]
		if t.compare(row[columnIndex], operator, value) {
			newSelectedRowIndices = append(newSelectedRowIndices, rowIndex)
		}
	}
	t.selectedRowIndices = newSelectedRowIndices
	return t
}

func (t *table) And(column, operator, value string) *table {
	return t.Where(column, operator, value)
}

func (t *table) Or(column, operator, value string) *table {
	columnIndex := t.columnToIndex[column]
	for rowIndex, row := range t.content {
		if t.compare(row[columnIndex], operator, value) {
			t.selectedRowIndices = append(t.selectedRowIndices, rowIndex)
		}
	}
	t.selectedRowIndices = lo.Uniq[int](t.selectedRowIndices)
	sort.Ints(t.selectedRowIndices)
	return t
}

func (t *table) Not(column, operator, value string) *table {
	newSelectedRowIndices := []int{}
	columnIndex := t.columnToIndex[column]
	for _, rowIndex := range t.selectedRowIndices {
		row := t.content[rowIndex]
		if !t.compare(row[columnIndex], operator, value) {
			newSelectedRowIndices = append(newSelectedRowIndices, rowIndex)
		}
	}
	t.selectedRowIndices = newSelectedRowIndices
	return t
}

func (t *table) compare(a, op, b string) bool {
	match := false
	switch op {
	case "=":
		match = a == b
	case "<>":
		match = a != b
	case "~~":
		match = regexp.MustCompilePOSIX(b).MatchString(a)
	case "<", ">", "<=", ">=":
		bInFloat, err := strconv.ParseFloat(b, 64)
		if err != nil {
			errors.Join(t.errs, err)
			return false // TODO: should exit early or bubble up the error earlier
		}
		aInFloat, err := strconv.ParseFloat(a, 64)
		if err != nil {
			errors.Join(t.errs, err)
			return false // TODO: should exit early or bubble up the error earlier
		}
		switch op {
		case "<":
			match = aInFloat < bInFloat
		case "<=":
			match = aInFloat <= bInFloat
		case ">":
			match = aInFloat > bInFloat
		case ">=":
			match = aInFloat >= bInFloat
		}
	}
	return match
}

func (t *table) Result() [][]string {
	if len(t.selectedColumns) == 0 {
		return t.content // TODO: deep copy return value
	}

	if len(t.selectedColumns) == 1 && t.selectedColumns[0] == "*" {
		return t.content // TODO: deep copy return value
	}

	indicesOfColumnsToKeep := make([]int, len(t.selectedColumns))
	for i, name := range t.selectedColumns {
		indicesOfColumnsToKeep[i] = t.columnToIndex[name]
	}
	sort.Ints(indicesOfColumnsToKeep)

	newContent := make([][]string, len(t.content))
	for i, row := range t.content {
		newContent[i] = make([]string, len(t.selectedColumns))
		for j, index := range indicesOfColumnsToKeep {
			newContent[i][j] = row[index]
		}
	}

	return newContent
}

func (t *table) Errors() error {
	return t.errs
}
