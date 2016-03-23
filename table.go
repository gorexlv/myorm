package orm

import (
	"database/sql"
	"fmt"
	"reflect"
)

type TableSchema struct {
	TABLE_CATALOG   string
	TABLE_SCHEMA    string
	TABLE_NAME      string
	TABLE_TYPE      string
	ENGINE          string
	VERSION         int
	ROW_FORMAT      string
	TABLE_ROWS      sql.NullInt64
	AVG_ROW_LENGTH  sql.NullInt64
	DATA_LENGTH     sql.NullInt64
	MAX_DATA_LENGTH sql.NullInt64
	INDEX_LENGTH    sql.NullInt64
	DATA_FREE       sql.NullInt64
	AUTO_INCREMENT  sql.NullInt64
	CREATE_TIME     sql.NullString
	UPDATE_TIME     sql.NullString
	CHECK_TIME      sql.NullString
	TABLE_COLLATION sql.NullString
	CHECKSUM        sql.NullInt64
	CREATE_OPTIONS  sql.NullString
	TABLE_COMMENT   sql.NullString
}

// Table table
type Table struct {
	orm     *ORM
	name    string
	charset string
	engine  string
	fields  map[string]ColField
	schemas map[string]ColSchema
	model   reflect.Type
}

func (table *Table) String() string {
	return fmt.Sprintf("%s : charset(%s) engine(%s) cols[%s]\n", table.name, table.charset, table.engine, table.fields)
}

func (table *Table) checkAlter() (changes, adds, drops []string) {
	for _, field := range table.fields {
		if _, ok := table.schemas[field.colName]; ok {
			changes = append(changes, field.colName)
		} else {
			adds = append(adds, field.colName)
		}
	}

	for _, schema := range table.schemas {
		if _, ok := table.fields[schema.COLUMN_NAME]; !ok {
			drops = append(drops, schema.COLUMN_NAME)
		}
	}

	return
}

// Alter alter
func (table *Table) Alter() (err error) {
	fmt.Println("table.Alter...")

	as := NewAlterStatement(table)
	as.Exec()

	return as.Error()
}

// Create create
func (table *Table) Create() (err error) {
	fmt.Println("table.Create...")
	return
}

func (table *Table) querySchemaInfo() ([]ColSchema, error) {
	sqlstr := "SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE table_schema = ? AND table_name = ?"
	var err error
	var rows *sql.Rows
	rows, err = table.orm.Query(sqlstr, table.orm.dbname, table.name)
	if err != nil {
		panic("fetch schema information error")
	}

	var schemas = make([]ColSchema, 0)
	for rows.Next() {
		var schema ColSchema
		if err := rows.Scan(&schema.TABLE_CATALOG, &schema.TABLE_SCHEMA, &schema.TABLE_NAME, &schema.COLUMN_NAME, &schema.ORDINAL_POSITION, &schema.COLUMN_DEFAULT, &schema.IS_NULLABLE, &schema.DATA_TYPE, &schema.CHARACTER_MAXIMUM_LENGTH, &schema.CHARACTER_OCTET_LENGTH, &schema.NUMERIC_PRECISION, &schema.NUMERIC_SCALE, &schema.DATETIME_PRECISION, &schema.CHARACTER_SET_NAME, &schema.COLLATION_NAME, &schema.COLUMN_TYPE, &schema.COLUMN_KEY, &schema.EXTRA, &schema.PRIVILEGES, &schema.COLUMN_COMMENT); err != nil {
			return schemas, err
		}

		schemas = append(schemas, schema)
	}

	return schemas, nil
}

func (table *Table) refreshSchemaInfo() error {
	schemas, err := table.querySchemaInfo()
	if err != nil {
		return err
	}

	for _, schema := range schemas {
		table.schemas[schema.COLUMN_NAME] = schema
	}

	return nil
}
