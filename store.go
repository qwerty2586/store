package store

import (
	"database/sql"
	"encoding/json"
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"reflect"
	"strings"
)

type Store struct {
	db         *sql.DB
	table_name string
}

func New(_db *sql.DB, _table_name string) (store *Store, err error) {
	//col := &xormColumn{tableName: _table_name}
	_, err = _db.Exec("CREATE TABLE IF NOT EXISTS " + _table_name + " (" +
		"`key` varchar(128) not null, " +
		"`value` text not null, " +
		"primary key (`key`) on conflict replace)")

	if err != nil {
		return
	}

	store = &Store{
		db:         _db,
		table_name: _table_name,
	}
	return
}

type Named interface {
	GetName() string
}

func (k *Store) Get(beans ...any) (err error) {
	defer err2.Return(&err)

	bean_names := getKeyNames(beans...)

	stmt := try.To1(k.db.Prepare("select * from " + k.table_name + " where `key` in (?" + strings.Repeat(",?", len(bean_names)-1) + ")"))
	defer stmt.Close()
	query := try.To1(stmt.Query(bean_names...))
	defer query.Close()

	for query.Next() {
		var key, value string
		try.To(query.Scan(&key, &value))
		for i, _ := range bean_names {
			if bean_names[i] == key {
				try.To(json.Unmarshal([]byte(value), beans[i]))
				break
			}
		}
	}
	return
}

func (k *Store) Set(beans ...any) (err error) {
	defer err2.Return(&err)
	bean_names := getKeyNames(beans...)

	s := "insert into " + k.table_name + " (`key`, `value`) values "
	vals := []interface{}{}

	for i, _ := range bean_names {
		s += "(?,?)"
		if i < len(bean_names)-1 {
			s += ","
		}
		vals = append(vals, bean_names[i], string(try.To1(json.Marshal(beans[i]))))
	}

	stmt := try.To1(k.db.Prepare(s))
	defer stmt.Close()
	try.To1(stmt.Exec(vals...))
	return
}

func getKeyNames(beans ...any) (keys []any) {
	named_inter := reflect.TypeOf((*Named)(nil)).Elem()

	keys = make([]any, len(beans))
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
