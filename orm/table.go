package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/oldbai555/lb/log"
	"reflect"
	"strings"
)

const (
	linePrefix = "  "
)

func (s *Session) Model(value interface{}) *Session {
	// nil or different model, update refTable
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *Schema {
	if s.refTable == nil {
		log.Errorf("Model is not set")
	}
	return s.refTable
}

func createTable(s *Session) error {
	table := s.RefTable()

	if len(table.Fields) == 0 {
		return errors.New("invalid field list")
	}

	primaryKeyCnt := 0

	// 查找主键
	// for _, f := range table.Fields {
	// 	rules := parseDbDef(f.Tag)
	//
	// 	if _, ok := rules.ruleMap["primary_key"]; ok {
	// 		primaryKeyCnt++
	// 		f.PrimaryKey = true
	// 		// int 的 primary key 加上自增
	// 		f.Extra = "AUTO_INCREMENT"
	// 	}
	//
	// }

	// 如果没有主键，拉 field_name = id 的来做主键
	if primaryKeyCnt == 0 {
		for _, f := range table.Fields {
			if f.Name == "id" && isIntType(f.Type) {
				primaryKeyCnt++
				f.PrimaryKey = true
				// int 的 primary key 加上自增
				f.Extra = "AUTO_INCREMENT"
			}
		}
	}
	// 主键最多一个
	if primaryKeyCnt > 1 || primaryKeyCnt == 0 {
		return fmt.Errorf("primary key count %d exceeded 1", primaryKeyCnt)
	}

	// 开始拼接创建表的语句
	var columns []string
	for _, field := range table.Fields {
		var items []string
		items = append(items, quoteName(field.Name), field.Type)

		l := strings.ToLower(field.Type)
		if field.PrimaryKey {
			items = append(items, "NOT NULL")
		} else if inSliceStr(l, []string{"text", "blob", "geometry", "json", "mediumtext"}) {
			items = append(items, "")
		} else if field.DefaultVal != "" {
			items = append(items, fmt.Sprintf("NOT NULL DEFAULT %s", field.DefaultVal))
		}

		if field.Extra != "" {
			items = append(items, field.Extra)
		}

		if field.Comment != "" {
			items = append(items, fmt.Sprintf("COMMENT '%s'", field.Comment))
		}

		createStmt := strings.Join(items, " ")
		columns = append(columns, linePrefix+createStmt)
	}

	for _, field := range table.Fields {
		if field.PrimaryKey {
			columns = append(columns, linePrefix+fmt.Sprintf("PRIMARY KEY (%s)", quoteName(field.Name)))
		}
	}

	desc := strings.Join(columns, ",\n")

	extra := "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"
	genSql := fmt.Sprintf("CREATE TABLE %s(\n%s\n)%s;", quoteName(table.Name), desc, extra)
	_, err := s.Raw(genSql).Exec()
	return err
}

func dropTable(s *Session) error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

func doDescTable(s *Session) (*descTable, error) {

	existSql, values := s.dialect.TableExistSQL(s.RefTable().Name)

	rows, err := s.Raw(existSql, values...).QueryRows()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	d := &descTable{}
	var field, typ, null, key, def, extra sql.NullString

	for rows.Next() {
		err = rows.Scan(&field, &typ, &null, &key, &def, &extra)
		if err != nil {
			log.Errorf("err:%s", err)
			return nil, err
		}
		d.columns = append(d.columns, &descTableColumn{
			Field:   field.String,
			Type:    typ.String,
			Null:    null.String,
			Key:     key.String,
			Default: def.String,
			Extra:   extra.String,
		})
	}
	return d, nil
}
