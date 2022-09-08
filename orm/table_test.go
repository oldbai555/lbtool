package orm

import (
	"github.com/jmoiron/sqlx"
	"github.com/oldbai555/lb/log"
	"reflect"
	"strings"
	"testing"
)

type User struct {
	Id   uint64
	Name string `lborm:"primary_key"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	engine, err := NewEngine(DMYSQL, "root:123456@tcp(175.178.156.14:3309)/orm")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	s := engine.NewSession().Model(&User{})
	err = s.DropTable()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	err = s.CreateTable()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
	return
}

func TestSession_DropTable(t *testing.T) {
	type fields struct {
		db       *sqlx.DB
		sql      strings.Builder
		sqlVars  []interface{}
		dialect  Dialect
		refTable *Schema
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				db:       tt.fields.db,
				sql:      tt.fields.sql,
				sqlVars:  tt.fields.sqlVars,
				dialect:  tt.fields.dialect,
				refTable: tt.fields.refTable,
			}
			if err := s.DropTable(); (err != nil) != tt.wantErr {
				t.Errorf("DropTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSession_HasTable(t *testing.T) {
	type fields struct {
		db       *sqlx.DB
		sql      strings.Builder
		sqlVars  []interface{}
		dialect  Dialect
		refTable *Schema
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				db:       tt.fields.db,
				sql:      tt.fields.sql,
				sqlVars:  tt.fields.sqlVars,
				dialect:  tt.fields.dialect,
				refTable: tt.fields.refTable,
			}
			if got := s.HasTable(); got != tt.want {
				t.Errorf("HasTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_Model(t *testing.T) {
	type fields struct {
		db       *sqlx.DB
		sql      strings.Builder
		sqlVars  []interface{}
		dialect  Dialect
		refTable *Schema
	}
	type args struct {
		value interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				db:       tt.fields.db,
				sql:      tt.fields.sql,
				sqlVars:  tt.fields.sqlVars,
				dialect:  tt.fields.dialect,
				refTable: tt.fields.refTable,
			}
			if got := s.Model(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSession_RefTable(t *testing.T) {
	type fields struct {
		db       *sqlx.DB
		sql      strings.Builder
		sqlVars  []interface{}
		dialect  Dialect
		refTable *Schema
	}
	tests := []struct {
		name   string
		fields fields
		want   *Schema
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				db:       tt.fields.db,
				sql:      tt.fields.sql,
				sqlVars:  tt.fields.sqlVars,
				dialect:  tt.fields.dialect,
				refTable: tt.fields.refTable,
			}
			if got := s.RefTable(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
