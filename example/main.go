package main

import (
	"database/sql"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/qwerty2586/store"
	"log"
)
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
	err2.Catch(func(err error) {
		log.Fatalln(err)
	})
	sq := try.To1(sql.Open("sqlite3", ":memory:"))
	kv := try.To1(store.New(sq, "test"))
	user := User{
		Id:   1,
		Name: "test_user1",
	}
	simple_int := SimpleInt(40)
	simple_string := SimpleString("test_simp2")
	try.To(kv.Set(&user, &simple_int, &simple_string))

	user = User{}
	simple_int = SimpleInt(0)
	simple_string = SimpleString("")

	try.To(kv.Get(&user, &simple_int, &simple_string))

	println(user.Id)
	println(user.Name)
	println(simple_int)
	println(simple_string)
}
