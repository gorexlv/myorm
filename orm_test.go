package orm

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name string
	Age  int
}

type UserBean struct {
	Bean  `name:"users" charset:"utf8" engine:"InnoDB"`
	Name  string `orm:"name varchar(30) null default 'foo'"`
	Age   string `orm:"age int(13) not null default 0"`
	Model *User
}

func Test_ORM(t *testing.T) {
	orm := New("root:abcdef1@tcp(localhost:3306)", "gorm", "utf8")
	orm.Modal(new(UserBean))
	for _, table := range orm.tables {
		fmt.Println(table)
	}
	if err := orm.AutoMigrate(); err != nil {
		panic(err)
	}
}
