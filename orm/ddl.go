package orm

import (
	"database/sql"
	"fmt"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/pkg/result"
	"strings"
)

// doDescTable 判断表是否存在,存在则返回表的字段
func doDescTable(s *Session) (*dbTable, error) {

	existSql, values := s.dialect.TableExistSQL(s.table.Name)

	rows, err := s.Raw(existSql, values...).QueryRows()
	if err != nil {
		if strings.Contains(err.Error(), "doesn't exist") &&
			strings.Contains(err.Error(), "1146: Table ") {
			return nil, result.NewLbErr(result.ErrOrmTableNotExist, err.Error())
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	d := &dbTable{}
	var field, typ, null, key, def, extra sql.NullString

	for rows.Next() {
		err = rows.Scan(&field, &typ, &null, &key, &def, &extra)
		if err != nil {
			log.Errorf("err:%s", err)
			return nil, err
		}
		d.columns = append(d.columns, &dbTableColumn{
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

// dropTable 删除表
func dropTable(s *Session) error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.table.Name)).Exec()
	return err
}

// genCreateTableSql 构建创表语句
func genCreateTableSql(s *Session) error {
	table := s.table

	if len(table.Fields) == 0 {
		return result.NewInvalidArg("invalid field list")
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

		field.createStmt = strings.Join(items, " ")
		columns = append(columns, linePrefix+field.createStmt)
	}

	for _, field := range table.Fields {
		if field.PrimaryKey {
			columns = append(columns, linePrefix+fmt.Sprintf("PRIMARY KEY (%s)", quoteName(field.Name)))
		}
	}

	desc := strings.Join(columns, ",\n")

	extra := "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin"
	table.createStmt = fmt.Sprintf("CREATE TABLE %s(\n%s\n)%s;", quoteName(table.Name), desc, extra)

	// _, err := s.Raw(table.createStmt).Exec()
	return nil
}

// createTable 创建表
func createTable(s *Session) error {
	err := genCreateTableSql(s)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	_, err = s.Raw(s.table.createStmt).Exec()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

// createOrUpdateTable 创建或更新表
func createOrUpdateTable(s *Session) error {
	table, err := doDescTable(s)

	if err != nil && result.GetErrCode(err) != result.ErrOrmTableNotExist {
		log.Errorf("err:%v", err)
		return err
	} else if err != nil && result.GetErrCode(err) == result.ErrOrmTableNotExist {
		// 找不到 那就创表
		err = createTable(s)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	// 构建创表语句
	err = genCreateTableSql(s)

	// 更新字段
	err = modifyTableColumn(s, table)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

// modifyTableColumn 修改表字段
func modifyTableColumn(s *Session, table *dbTable) error {

	// 比一下要加的新列
	fieldMap2Create := map[string]bool{}
	var addColumns []string

	refTable := s.table

	for _, f := range refTable.Fields {
		found := false
		for _, x := range table.columns {
			if strings.EqualFold(f.Name, x.Field) {
				found = true
				break
			}
		}
		if !found {
			fieldMap2Create[f.Name] = true
			addColumns = append(addColumns, f.Name)
		}
	}

	// 新增字段
	if len(fieldMap2Create) > 0 {
		for _, f := range refTable.Fields {
			if !fieldMap2Create[f.Name] {
				continue
			}
			stmt := fmt.Sprintf("ALTER TABLE %s ADD %s", quoteName(refTable.Name), f.createStmt)
			_, err := s.Raw(stmt).Exec()
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		}
	}

	return nil
}
