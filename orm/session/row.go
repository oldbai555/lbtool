package session

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/orm/dialect"
	"github.com/oldbai555/lb/orm/schema"
	"reflect"
	"strings"
)

type Session struct {
	// db 连接数据库成功的连接
	db *sqlx.DB

	// sql 拼接 SQL 语句
	sql strings.Builder
	// sqlVars SQL 语句中占位符的对应值
	sqlVars []interface{}

	// Dial 数据类型转换
	Dial dialect.Dialect
	// Table 映射表
	Table *schema.Schema
}

func NewSession(db *sqlx.DB, dial dialect.Dialect) *Session {
	return &Session{
		db:   db,
		Dial: dial,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sqlx.DB {
	return s.db
}

// Raw 填充 SQL 中的占位符
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec raw sql with sqlVars
func (s *Session) Exec() (result sql.Result, err error) {
	defer func() {
		s.Clear()
	}()
	log.Infof(strings.ReplaceAll(s.sql.String(), "?", "%v"), s.sqlVars...)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Errorf("err:%v", err)
	}
	return
}

// QueryRow gets a record from db
func (s *Session) QueryRow() *sql.Row {
	log.Infof(strings.ReplaceAll(s.sql.String(), "?", "%v"), s.sqlVars...)
	defer func() {
		s.Clear()
	}()
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows gets a list of records from db
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Infof(strings.ReplaceAll(s.sql.String(), "?", "%v"), s.sqlVars...)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Errorf("err:%v", err)
	}
	return
}

func (s *Session) Model(value interface{}) *Session {
	// nil or different model, update refTable
	if s.Table == nil ||
		reflect.TypeOf(value) != reflect.TypeOf(s.Table.Model) {
		s.Table = schema.Parse(value, s.Dial)
	}
	return s
}
