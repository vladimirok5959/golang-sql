package common

import (
	"database/sql"
	"reflect"
)

type Rows struct {
	*sql.Rows
}

func scans(row any) []any {
	v := reflect.ValueOf(row).Elem()
	res := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		res[i] = v.Field(i).Addr().Interface()
	}
	return res
}

func (r *Rows) Scans(row any) error {
	return r.Rows.Scan(scans(row)...)
}
