package store

import (
	"encoding/json"
	"reflect"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type Store struct {
	e          *xorm.Engine
	table_name string
}

type xormColumn struct {
	Key       string `xorm:"varchar(128) notnull pk"`
	Value     string `xorm:"text notnull"`
	tableName string `xorm:"ignore"`
}

func (xc *xormColumn) TableName() string {
	return xc.tableName
}

func New(_e *xorm.Engine, _table_name string) *Store {
	//col := &xormColumn{tableName: _table_name}
	_e.Exec("CREATE TABLE IF NOT EXISTS " + _table_name + " (" +
		"`key` varchar(128) not null, " +
		"`value` text not null, " +
		"primary key (`key`) on conflict replace)")
	return &Store{e: _e, table_name: _table_name}
}

type Named interface {
	GetName() string
}

func (k *Store) Get(beans ...any) (err error) {
	bean_names := getKeyNames(beans...)

	var cols []xormColumn
	err = k.e.Table(k.table_name).Where(builder.In("key", bean_names)).Find(&cols)
	if err != nil {
		return
	}

	for _, col := range cols {
		key := col.Key
		for i, _ := range bean_names {
			if bean_names[i] == key {
				err = json.Unmarshal([]byte(col.Value), beans[i])
				if err != nil {
					return
				}
				break
			}
		}
	}
	return
}

func (k *Store) Set(beans ...any) (err error) {
	bean_names := getKeyNames(beans...)

	cols := make([]xormColumn, len(beans))
	for i, bean := range beans {
		cols[i].Key = bean_names[i]
		var b []byte
		b, err = json.Marshal(bean)
		if err != nil {
			return
		}
		cols[i].Value = string(b)
	}
	_, err = k.e.Table(k.table_name).Insert(cols)
	return
}

func getKeyNames(beans ...any) (keys []string) {
	named_inter := reflect.TypeOf((*Named)(nil)).Elem()

	keys = make([]string, len(beans))
	for i, bean := range beans {
		bean_type := reflect.TypeOf(bean)
		if bean_type.Kind() != reflect.Ptr {
			panic("bean must be a pointer")
		}
		elem := bean_type.Elem()
		keys[i] = elem.Name()
		if bean_type.Implements(named_inter) {
			keys[i] = bean.(Named).GetName()
		}
	}
	return
}
