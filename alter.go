package orm

import (
	"errors"
	"fmt"
	"strings"
)

// AlterStatement alter statment
type AlterStatement struct {
	table    *Table
	changes  []string
	addes    []string
	drops    []string
	indexes  []string
	primaris []string
	uniques  []string
	sql      string
	err      error
}

// NewAlterStatement returns new AlterStatement
func NewAlterStatement(table *Table) *AlterStatement {
	return &AlterStatement{
		table: table,
	}
}

func (as *AlterStatement) explain() {
	if as.err != nil {
		return
	}

	sql := "ALTER TABLE " + ("`" + as.table.name + "`")

	var items []string

	for _, item := range as.addes {
		f := as.table.fields[item]
		if f.colName != "" {
			items = append(items, f.addStmt())
		}
	}
	for _, item := range as.changes {
		f := as.table.fields[item]
		if f.colName != "" {
			items = append(items, " MODIFY "+"`"+f.colName+"` "+f.changeStmt())
		}
	}
	for _, item := range as.drops {
		f := as.table.fields[item]
		if f.colName != "" {
			items = append(items, " DROP "+" `"+f.colName+"` "+f.dropStmt())
		}
	}

	sql += strings.Join(items, ",")

	fmt.Println("stmt.sql: ", sql)
	as.sql = sql
}

func (as *AlterStatement) check() {
	if as.err != nil {
		return
	}

	schemas, err := as.table.querySchemaInfo()
	if err != nil {
		as.err = err
		return
	}

	if len(schemas) == 0 {
		as.err = errors.New("query column schema failed")
		return
	}

	for _, field := range as.table.fields {
		var ok bool
		for _, schema := range schemas {
			if schema.COLUMN_NAME == field.colName {
				ok = true
			}
		}
		if ok {
			as.changes = append(as.changes, field.colName)
		} else {
			as.addes = append(as.addes, field.colName)
		}
	}

	for _, schema := range schemas {
		if _, ok := as.table.fields[schema.COLUMN_NAME]; !ok {
			as.drops = append(as.drops, schema.COLUMN_NAME)
		}
	}
}

func (as *AlterStatement) Error() error {
	return as.err
}

func (as *AlterStatement) Exec() {
	as.check()
	as.explain()
	_, err := as.table.orm.Exec(as.sql)
	if err != nil {
		as.err = err
		panic(err)
	}
}
