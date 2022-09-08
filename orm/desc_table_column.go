package orm

type descTableColumn struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

type descTable struct {
	columns []*descTableColumn
}
