package orm

// Statement statement
type Statement interface {
	Exec() // 只执行不返回结果
	Query()
}

// SelectStatement select statement
type SelectStatement struct {
	tables map[string]*Table
}

// NewSelectStatement returns new SelectStatement
func NewSelectStatement(tables ...Table) *SelectStatement {
	return &SelectStatement{}
}
