package orm

import (
	"fmt"
	"reflect"
)

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// DataTypeOf 用于将 Go 语言的类型转换为该数据库的数据类型
	DataTypeOf(typ reflect.Value) string
	// TableExistSQL 返回某个表是否存在的 SQL 语句，参数是表名(table)
	TableExistSQL(tableName string) (string, []interface{})
	// GetFieldDefaultValue 得到默认值
	GetFieldDefaultValue(typ reflect.Value) string
}

// RegisterDialect 注册 dialect 实例
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// getDialect 获取 dialect 实例
func getDialect(name string) (dialect Dialect, err error) {
	var ok bool
	dialect, ok = dialectsMap[name]
	if !ok {
		err = fmt.Errorf("not found dialect,name: %s", name)
	}
	return
}
