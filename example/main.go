package main

import (
	"github.com/qwerty2586/store"
	"xorm.io/xorm"
)
import _ "github.com/go-sql-driver/mysql"
import _ "github.com/mattn/go-sqlite3"

type User struct {
	Id   int64
	Name string
}

type SimpleInt int

func (i *SimpleInt) GetName() string {
	return "simple_int_with_special_name"
}

type SimpleString string

func (i *SimpleString) GetName() string {
	return "simple_string_with_special_name"
}

func main() {
	e, _ := xorm.NewEngine("sqlite3", "./test.db")
	e.ShowSQL(true)
	kv := store.New(e, "test")
	user := User{
		Id:   1,
		Name: "test_user",
	}
	simple_int := SimpleInt(40)
	simple_string := SimpleString("test_simp")
	kv.Set(&user, &simple_int, &simple_string)

	user = User{}
	simple_int = SimpleInt(0)
	simple_string = SimpleString("")

	kv.Get(&user, &simple_int, &simple_string)

	println(user.Id)
	println(user.Name)
	println(simple_int)
	println(simple_string)
}
