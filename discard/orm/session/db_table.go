package session

// DbTableColumn 数据库的表字段
type DbTableColumn struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

// DbTable 数据库的表信息
type DbTable struct {
	Columns []*DbTableColumn
}
