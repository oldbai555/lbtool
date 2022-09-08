package orm

import (
	"go/ast"
	"reflect"
	"strings"
)

// Field represents a column of database
type Field struct {
	Name       string
	Type       string
	Tag        string
	DefaultVal string
	PrimaryKey bool
	Comment    string
	Extra      string
	createStmt string
}

// Schema represents a table of database
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
	createStmt string
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// Parse 解析结构体转换为数据库表
func Parse(dest interface{}, d Dialect) *Schema {
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
				Name:       strings.ToLower(p.Name),
				Type:       d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
				DefaultVal: d.GetFieldDefaultValue(reflect.Indirect(reflect.New(p.Type))),
				Comment:    strings.ToLower(p.Name),
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
