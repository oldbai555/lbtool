package schema

import (
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/orm/dialect"
	"go/ast"
	"reflect"
	"strings"
)

// Field represents a column of database
type Field struct {
	DbName     string
	Type       string
	Tag        string
	DefaultVal string
	PrimaryKey bool
	Comment    string
	Extra      string
	CreateStmt string

	// OriginalName 原始名字
	OriginalName string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
	CreateStmt string
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldByName := destValue.FieldByName(field.OriginalName)
		if !fieldByName.IsValid() {
			log.Warnf("field name %s is valid, field by %s", field.OriginalName, field.OriginalName)
			continue
		}
		fieldValues = append(fieldValues, fieldByName.Interface())
	}
	return fieldValues
}

// Parse 解析结构体转换为数据库表
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     strings.ToLower(modelType.Name()),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				DbName:       strings.ToLower(p.Name),
				Type:         d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
				DefaultVal:   d.GetFieldDefaultValue(reflect.Indirect(reflect.New(p.Type))),
				Comment:      p.Name,
				OriginalName: p.Name,
			}
			if v, ok := p.Tag.Lookup("lborm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
