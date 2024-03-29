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
	db          *sql.DB
	table_name  string
	driver_name string
}

func New(_driver_name string, _db *sql.DB, _table_name string) (store *Store, err error) {

	if _driver_name != "mysql" && _driver_name != "sqlite3" {
		panic("driver " + _driver_name + " not supported")
	}

	// sanitize table name, put out chars that are not allowed in table names
	_table_name = sanitizeTableName(_table_name)

	_, err = _db.Exec("CREATE TABLE IF NOT EXISTS " + _table_name + " (" +
		"`kv_key` varchar(128) not null, " +
		"`kv_value` text not null, " +
		"primary key (`kv_key`))")

	if err != nil {
		return
	}

	store = &Store{
		db:          _db,
		table_name:  _table_name,
		driver_name: _driver_name,
	}
	return
}

// alow only alphanumeric characters and underscores
func sanitizeTableName(name string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' {
			return r
		}
		return -1
	}, name)
}

type Named interface {
	GetName() string
}

func (k *Store) Get(beans ...any) (found int64, err error) {
	defer err2.Return(&err)

	bean_names := getKeyNames(beans...)

	stmt := try.To1(k.db.Prepare("select * from " + k.table_name + " where `kv_key` in (?" + strings.Repeat(",?", len(bean_names)-1) + ")"))
	defer stmt.Close()
	query := try.To1(stmt.Query(bean_names...))
	defer query.Close()

	for query.Next() {
		var key, value string
		try.To(query.Scan(&key, &value))
		for i, _ := range bean_names {
			if bean_names[i] == key {
				try.To(json.Unmarshal([]byte(value), beans[i]))
				found++
				// faster reslicing
				bean_names[i] = bean_names[len(bean_names)-1]
				bean_names = bean_names[:len(bean_names)-1]
				beans[i] = beans[len(beans)-1]
				beans = beans[:len(beans)-1]
				break
			}
		}
	}
	// use reflec.Zero to set the rest of the beans to zero
	for i, bean_name := range bean_names {
		if bean_name != "" {
			reflect.ValueOf(beans[i]).Elem().Set(reflect.Zero(reflect.ValueOf(beans[i]).Elem().Type()))
		}
	}
	return
}

func (k *Store) Set(beans ...any) (err error) {
	defer err2.Return(&err)
	bean_names := getKeyNames(beans...)

	s := "insert into " + k.table_name + " (`kv_key`, `kv_value`) values "
	vals := []interface{}{}

	for i, _ := range bean_names {
		s += "(?,?)"
		if i < len(bean_names)-1 {
			s += ", "
		}
		str_value := string(try.To1(json.Marshal(beans[i])))
		vals = append(vals, bean_names[i], str_value)
	}
	if k.driver_name == "mysql" {
		s += "ON DUPLICATE KEY UPDATE kv_value=VALUES(kv_value)"
	}
	if k.driver_name == "sqlite3" {
		s += "ON CONFLICT(kv_key) DO UPDATE SET kv_value=excluded.kv_value"
	}

	stmt := try.To1(k.db.Prepare(s))
	defer stmt.Close()
	try.To1(stmt.Exec(vals...))
	return
}

func (k *Store) Delete(beans ...any) (deleted int64, err error) {
	defer err2.Return(&err)
	bean_names := getKeyNames(beans...)

	s := "delete from " + k.table_name + " where `kv_key` in (?" + strings.Repeat(",?", len(bean_names)-1) + ")"
	stmt := try.To1(k.db.Prepare(s))
	defer stmt.Close()
	res := try.To1(stmt.Exec(bean_names...))
	deleted, _ = res.RowsAffected()
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
