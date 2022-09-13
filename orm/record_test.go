package orm

import (
	"github.com/oldbai555/lb/log"
	"github.com/oldbai555/lb/orm/dialect"
	"github.com/oldbai555/lb/orm/session"
	"testing"
)

var (
	user1 = &User{
		Name: "Tom",
		Age:  18,
	}
	user2 = &User{
		Name: "Sam",
		Age:  25,
	}
	user3 = &User{
		Name: "Jack",
		Age:  25,
	}
)

func testRecordInit(t *testing.T) *session.Session {
	t.Helper()
	engine, err := NewEngine(dialect.DMYSQL, "root:123456@tcp(175.178.156.14:3309)/orm")
	if err != nil {
		log.Errorf("err:%v", err)
		return nil
	}
	s := engine.NewSession().Model(&User{})
	err1 := dropTable(s)
	err2 := createTable(s)
	_, err3 := s.Insert(user1, user2)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatal("failed init test records")
	}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
}
