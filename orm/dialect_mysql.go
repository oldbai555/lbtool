package orm

import (
	"fmt"
	"reflect"
	"time"
)

const (
	DMYSQL = "mysql"
)

func init() {
	RegisterDialect(DMYSQL, &dMysql{})
}

type dMysql struct{}

func (m *dMysql) GetFieldDefaultValue(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return "0"
	case reflect.Float32, reflect.Float64:
		return "0.0"
	default:
		return ""
	}
}

func (m *dMysql) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "TINYINT(1)"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "INT(10)"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "INT(10) UNSIGNED"
	case reflect.Int64:
		return "BIGINT(20)"
	case reflect.Uint64:
		return "BIGINT(20) UNSIGNED"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Array, reflect.Slice:
		return "BLOB"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(any(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind())))
}

func (m *dMysql) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "select TABLE_NAME from information_schema.TABLES where TABLE_NAME = '?';", args
}

var _ Dialect = (*dMysql)(nil)
