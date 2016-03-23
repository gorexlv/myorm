package orm

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

type ColSchema struct {
	TABLE_CATALOG            string
	TABLE_SCHEMA             string
	TABLE_NAME               string
	COLUMN_NAME              string
	ORDINAL_POSITION         int
	COLUMN_DEFAULT           sql.NullString
	IS_NULLABLE              string
	DATA_TYPE                string
	CHARACTER_MAXIMUM_LENGTH sql.NullInt64
	CHARACTER_OCTET_LENGTH   sql.NullInt64
	NUMERIC_PRECISION        sql.NullInt64
	NUMERIC_SCALE            sql.NullInt64
	DATETIME_PRECISION       sql.NullInt64
	CHARACTER_SET_NAME       sql.NullString
	COLLATION_NAME           sql.NullString
	COLUMN_TYPE              string
	COLUMN_KEY               string
	EXTRA                    string
	PRIVILEGES               string
	COLUMN_COMMENT           string
}

// Column column
type ColField struct {
	modelName  string
	colName    string
	colType    string
	refName    string
	refType    string
	notNull    bool
	isPK       bool
	isFK       bool
	isIndex    bool
	unique     bool
	autoIncr   bool
	comment    string
	defaultVal interface{}
	tag        string
}

func (col ColField) String() string {
	return fmt.Sprintf(" %s \n", col.tag)
}

func (col ColField) alter() string {
	var str = col.colType

	if col.notNull {
		str += " NOT NULL"
	} else {
		str += " NULL"
	}

	if col.autoIncr {
		str += " AUTO_INCREMENT"
	}

	if col.defaultVal != nil {
		var val string
		switch col.defaultVal.(type) {
		case string:
			val = "'" + col.defaultVal.(string) + "'"
		default:
			val = fmt.Sprint(val)
		}
		str += " DEFAULT " + val
	}

	return str
}

var (
	intexp     *regexp.Regexp
	floatexp   *regexp.Regexp
	doubleexp  *regexp.Regexp
	decimalexp *regexp.Regexp
	varcharexp *regexp.Regexp
)

func init() {
	intexp, _ = regexp.Compile(`(int|int\((\d+)\))`)
	floatexp, _ = regexp.Compile(`(float|float(\d+))`)
	doubleexp, _ = regexp.Compile(`(double|double(\d+))`)
	decimalexp, _ = regexp.Compile(`(decimal|decimal(\d+))`)
	varcharexp, _ = regexp.Compile(`varchar\((\d+)\)`)
}

// Parse parse
func (col ColField) Parse() ColField {
	if splits := strings.Split(col.tag, " "); len(splits) > 0 {
		col.colName = splits[0]
	} else {
		panic("invalid tag")
	}

	tag := strings.ToLower(col.tag)
	fmt.Println("attr: ", tag)
	if strings.Contains(tag, "int") {
		col.colType = "int"
	}

	if typ := intexp.FindString(tag); typ != "" {
		col.colType = typ
	}
	if typ := floatexp.FindString(tag); typ != "" {
		col.colType = typ
	}
	if typ := doubleexp.FindString(tag); typ != "" {
		col.colType = typ
	}
	if typ := decimalexp.FindString(tag); typ != "" {
		col.colType = typ
	}
	if typ := varcharexp.FindString(tag); typ != "" {
		col.colType = typ
	}

	if strings.Contains(tag, "text") {
		col.colType = "text"
	}
	if strings.Contains(tag, "datetime") {
		col.colType = "datetime"
	}
	if strings.Contains(tag, "date ") {
		col.colType = "date"
	}
	if strings.Contains(tag, "time ") {
		col.colType = "time"
	}
	if strings.Contains(tag, "timestamp") {
		col.colType = "timestamp"
	}
	if strings.Contains(tag, " auto") {
		col.autoIncr = true
	}
	if strings.Contains(tag, " not null ") {
		col.notNull = true
	}
	if strings.Contains(tag, "primary key") {
		col.isPK = true
	}

	fmt.Println("col:", col.colType, col.notNull, col.autoIncr)

	return col
}
