package orm

// dbTableColumn 数据库的表字段
type dbTableColumn struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

// dbTable 数据库的表信息
type dbTable struct {
	columns []*dbTableColumn
}
