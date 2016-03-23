package orm

import (
	"database/sql"
	"reflect"
)

// Bean bean
type Bean interface{}

// ORM orm
type ORM struct {
	*sql.DB
	dbname  string
	charset string
	tables  map[string]*Table
	schemas map[string]*TableSchema
}

// New return a new ORM instance
func New(addr, dbname, charset string) *ORM {
	db, err := sql.Open("mysql", addr+"/"+dbname+"?charset="+charset)
	if err != nil {
		panic(err)
	}

	return &ORM{
		DB:      db,
		dbname:  dbname,
		charset: charset,
		tables:  make(map[string]*Table),
	}
}

// Modal modal
func (orm *ORM) Modal(beans ...interface{}) {
	for _, bean := range beans {
		bt := reflect.TypeOf(bean).Elem()
		bv := reflect.ValueOf(bean).Elem()

		table := &Table{
			orm:     orm,
			fields:  make(map[string]ColField),
			schemas: make(map[string]ColSchema),
		}

		for index := 0; index < bt.NumField(); index++ {
			fType := bt.Field(index)
			fName := fType.Name
			fValue := bv.FieldByName(fName)

			switch {
			case reflect.String == fValue.Type().Kind():
				col := ColField{
					modelName: fName,
					tag:       fType.Tag.Get("orm"),
				}.Parse()

				// table.fields = append(table.fields, col)
				table.fields[col.colName] = col
			case "Model" == fName:
				table.model = fValue.Type()
			case "Bean" == fName:
				table.name = fType.Tag.Get("name")
				table.engine = fType.Tag.Get("engine")
				table.charset = fType.Tag.Get("charset")
			}
		}

		orm.tables[table.name] = table
	}
}

func (orm *ORM) querySchemaInfo() ([]TableSchema, error) {
	sqlstr := "SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ?"
	var err error
	var rows *sql.Rows
	rows, err = orm.Query(sqlstr, orm.dbname)
	if err != nil {
		panic("fetch schema information error")
	}

	var schemas = make([]TableSchema, 0)
	for rows.Next() {
		var schema TableSchema
		if err := rows.Scan(&schema.TABLE_CATALOG, &schema.TABLE_SCHEMA, &schema.TABLE_NAME, &schema.TABLE_TYPE, &schema.ENGINE, &schema.VERSION, &schema.ROW_FORMAT, &schema.TABLE_ROWS, &schema.AVG_ROW_LENGTH, &schema.DATA_LENGTH, &schema.MAX_DATA_LENGTH, &schema.INDEX_LENGTH, &schema.DATA_FREE, &schema.AUTO_INCREMENT, &schema.CREATE_TIME, &schema.UPDATE_TIME, &schema.CHECK_TIME, &schema.TABLE_COLLATION, &schema.CHECKSUM, &schema.CREATE_OPTIONS, &schema.TABLE_COMMENT); err != nil {
			return schemas, err
		}

		schemas = append(schemas, schema)
	}

	return schemas, nil
}

func (orm *ORM) RefreshSchemaInfo(recursive bool) (err error) {
	schemas, err := orm.querySchemaInfo()
	if err != nil {
		return
	}

	for _, schema := range schemas {
		orm.schemas[schema.TABLE_NAME] = &schema
	}

	if recursive {
		// 刷新table的schema信息
		for _, table := range orm.tables {
			table.refreshSchemaInfo()
		}
	}

	return nil
}

// AutoMigrate auto migrate
func (orm *ORM) AutoMigrate() error {
	schemas, err := orm.querySchemaInfo()
	if err != nil {
		return err
	}

	for _, table := range orm.tables {
		var ok bool
		for _, schema := range schemas {
			if schema.TABLE_NAME == table.name {
				ok = true
			}
		}

		if ok {
			err = table.Alter()
		} else {
			err = table.Create()
		}
	}

	return nil
}
