package orm

import (
	"reflect"
)

const (
	linePrefix = "  "
)

func (s *Session) Model(value interface{}) *Session {
	// nil or different model, update refTable
	if s.table == nil ||
		reflect.TypeOf(value) != reflect.TypeOf(s.table.Model) {
		s.table = Parse(value, s.dialect)
	}
	return s
}
