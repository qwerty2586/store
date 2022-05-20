package store

import (
	"database/sql"
	"github.com/lainio/err2"
	"github.com/lainio/err2/assert"
	"github.com/lainio/err2/try"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

var (
	testStore *Store
)

const test_table_name = "test_table"

func TestMain(m *testing.M) {
	defer err2.Catch(func(err error) {
		panic(err)
	})
	// sql memory db
	db := try.To1(sql.Open("sqlite3", ":memory:"))

	testStore = try.To1(New(db, test_table_name))
	m.Run()
	try.To(db.Close())
}

func TestInt(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom int
	type MyInt int

	in := MyInt(42)
	try.To(testStore.Set(&in))

	var out MyInt
	try.To1(testStore.Get(&out))

	assert.P.True(in == out, "returned value should be the same as set value")
}

func TestString(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom string
	type MyString string

	in := MyString("hello")
	try.To(testStore.Set(&in))

	var out MyString
	try.To1(testStore.Get(&out))

	assert.P.True(in == out, "returned value should be the same as set value")
}

func TestSlice(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom slice
	type MySlice []string

	in := MySlice([]string{"hello", "world"})
	try.To(testStore.Set(&in))

	var out MySlice
	try.To1(testStore.Get(&out))

	assert.P.True(len(in) == len(out), "returned slice length should be the same as set value")

	for i, v := range in {
		assert.P.True(v == out[i], "returned slice value should be the same as set value")
	}
}

func TestMap(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom map
	type MyMap map[string]string

	in := MyMap(map[string]string{
		"hello":        "world",
		"How are you?": "fine",
	})
	try.To(testStore.Set(&in))

	var out MyMap
	try.To1(testStore.Get(&out))

	assert.P.True(len(in) == len(out), "returned map length should be the same as set value")

	for k, v := range in {
		assert.P.True(v == out[k], "returned map value should be the same as set value")
	}
}

func TestStruct(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom struct
	type MyStruct struct {
		Name string
		Age  int
	}

	in := MyStruct{
		Name: "John",
		Age:  42,
	}
	try.To(testStore.Set(&in))

	var out MyStruct
	try.To1(testStore.Get(&out))

	assert.P.True(in.Name == out.Name, "returned struct value should be the same as set value")
	assert.P.True(in.Age == out.Age, "returned struct value should be the same as set value")
}

func TestMultiple(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom int
	type MyInt int
	in1 := MyInt(42)
	type MyString string
	in2 := MyString("hello2")
	type MySlice []string
	in3 := MySlice([]string{"hello2", "world2"})
	type MyMap map[string]string
	in4 := MyMap(map[string]string{
		"hello2":        "world2",
		"How are you?2": "fine2",
	})
	type MyStruct struct {
		Name string
		Age  int
	}
	in5 := MyStruct{
		Name: "John2",
		Age:  43,
	}
	try.To(testStore.Set(&in1, &in2, &in3, &in4, &in5))

	var out1 MyInt
	var out2 MyString
	var out3 MySlice
	var out4 MyMap
	var out5 MyStruct

	// out of order
	count := try.To1(testStore.Get(&out5, &out2, &out1, &out4, &out3))

	assert.P.True(in1 == out1, "returned value should be the same as set value")
	assert.P.True(in2 == out2, "returned value should be the same as set value")
	assert.P.True(len(in3) == len(out3), "returned slice length should be the same as set value")
	assert.P.True(len(in4) == len(out4), "returned map length should be the same as set value")
	assert.P.True(in5.Name == out5.Name, "returned struct value should be the same as set value")
	assert.P.True(in5.Age == out5.Age, "returned struct value should be the same as set value")
	assert.P.True(count == 5, "returned count should be 5")
}

func TestZeroingUnknownInt(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom int
	type MyIntX int
	out := MyIntX(42)
	try.To1(testStore.Get(&out))
	assert.P.True(out == 0, "unknown store value should be zeroed")
}

func TestZeroingUnknownString(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom int
	type MyStringX string
	out := MyStringX("hello")
	try.To1(testStore.Get(&out))
	assert.P.True(out == "", "unknown store value should be zeroed")
}

func TestZeroingUnknownSlice(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom int
	type MySliceX []string
	out := MySliceX([]string{"hello"})
	try.To1(testStore.Get(&out))
	assert.P.True(len(out) == 0, "unknown store value should be zeroed")
}

func TestZeroUnknownMap(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom int
	type MyMapX map[string]string
	out := MyMapX(map[string]string{
		"hello": "world",
	})
	count := try.To1(testStore.Get(&out))
	assert.P.True(len(out) == 0, "unknown store value should be zeroed")
	assert.P.True(count == 0, "returned count should be zero")
}

func TestDeletingUnknown(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom int
	type MyIntXX int
	out := MyIntXX(42)
	deleted := try.To1(testStore.Delete(&out))
	assert.P.True(deleted == 0, "deleted count should be zero")
}

func TestDeletingExisting(t *testing.T) {
	defer err2.Catch(func(err error) {
		t.Error(err)
	})
	// custom int
	type MyStructXX struct {
		Name string
		Age  int
	}
	in := MyStructXX{
		Name: "John",
		Age:  42,
	}
	try.To(testStore.Set(&in))
	deleted := try.To1(testStore.Delete(&MyStructXX{}))
	assert.P.True(deleted == 1, "deleted count should be 1")
	// now we try to get it
	// non zero values
	out := MyStructXX{
		Name: "Kate",
		Age:  43,
	}
	count := try.To1(testStore.Get(&out))
	assert.P.True(count == 0, "returned count should be zero")
	assert.P.True(out.Name == "", "returned struct value should be zeroed")
	assert.P.True(out.Age == 0, "returned struct value should be zeroed")
}
